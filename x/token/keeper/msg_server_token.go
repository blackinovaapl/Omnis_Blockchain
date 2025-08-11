

package keeper

import (
	"context"
	"errors" // Keep errors for collections.ErrNotFound
	"fmt"
	"strconv" // New: Needed for converting ID to string for events

	"omnis/x/token/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types" // New: Needed for coin operations and UnwrapSDKContext
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateToken(goCtx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx) // Unwrap to sdk.Context

	creatorAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if a token with the same symbol already exists (we'll implement GetTokenBySymbol in keeper.go)
	_, found := k.GetTokenBySymbol(ctx, msg.Symbol)
	if found {
		return nil, errorsmod.Wrapf(types.ErrTokenAlreadyExists, "token with symbol %s already exists", msg.Symbol)
	}

	// Convert totalSupply string to a proper sdk.Int for calculations
	totalSupplyInt, ok := sdk.NewIntFromString(msg.TotalSupply)
	if !ok {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid total supply: %s", msg.TotalSupply)
	}
	if totalSupplyInt.IsNegative() {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "total supply cannot be negative")
	}

	nextId, err := k.TokenSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	var token = types.Token{
		Id:          nextId,
		Creator:     msg.Creator,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Decimals:    msg.Decimals,
		TotalSupply: msg.TotalSupply, // Store as string
		Metadata:    msg.Metadata,
	}

	if err = k.Token.Set(
		ctx,
		nextId,
		token,
	); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set token")
	}

	// Mint the initial supply and send it to the creator
	// Define the coin for the new token. The denom will be the token symbol.
	if !sdk.IsValidDenom(msg.Symbol) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid token symbol for denom: %s", msg.Symbol)
	}
	coin := sdk.NewCoin(msg.Symbol, totalSupplyInt)
	coins := sdk.NewCoins(coin)

	// Mint coins to the module account
	// The module account name is types.ModuleName (which is "token")
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "failed to mint coins: %v", err)
	}

	// Send minted coins from module account to creator
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "failed to send minted coins to creator: %v", err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeCreateToken,
			sdk.NewAttribute(types.AttributeKeyTokenID, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyTokenName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyTokenSymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyTotalSupply, msg.TotalSupply),
		),
	)

	return &types.MsgCreateTokenResponse{
		Id: nextId,
	}, nil
}



func (k msgServer) UpdateToken(ctx context.Context, msg *types.MsgUpdateToken) (*types.MsgUpdateTokenResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	var token = types.Token{
		Creator:     msg.Creator,
		Id:          msg.Id,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Decimals:    msg.Decimals,
		TotalSupply: msg.TotalSupply,
		Metadata:    msg.Metadata,
	}

	// Checks that the element exists
	val, err := k.Token.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get token")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Token.Set(ctx, msg.Id, token); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update token")
	}

	return &types.MsgUpdateTokenResponse{}, nil
}

func (k msgServer) DeleteToken(ctx context.Context, msg *types.MsgDeleteToken) (*types.MsgDeleteTokenResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.Token.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get token")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Token.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete token")
	}

	return &types.MsgDeleteTokenResponse{}, nil
}
