
package keeper

import (
	"context" // New: Needed for collections interactions
	"fmt"

	"omnis/x/token/types"

	"cosmossdk.io/collections" // Keep collections
	"cosmossdk.io/core/address"
	"cosmossdk.io/log" // Use cosmossdk.io/log instead of tendermint/tendermint/libs/log

	sdk "github.com/cosmos/cosmos-sdk/types" // New: Needed for bankKeeper
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types" // New: Needed for accountKeeper if not already there
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper" // New: Import bankkeeper
)

type (
	Keeper struct {
		// Keep the collections that Ignite scaffolded
		Token    collections.Map[uint64, types.Token]
		TokenSeq collections.Sequence
		Schema   collections.Schema

		addressCodec address.Codec
		bankKeeper   bankkeeper.Keeper       // New: Add bankKeeper
		accountKeeper types.AccountKeeper    // New: Add accountKeeper (use types.AccountKeeper to match interface)
	}
)

func NewKeeper(
	addressCodec address.Codec,
	storeService sdk.StoreService, // Use StoreService from SDK context
	bankKeeper bankkeeper.Keeper,     // New: Pass bankKeeper
	accountKeeper types.AccountKeeper, // New: Pass accountKeeper
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		addressCodec: addressCodec,
		bankKeeper:   bankKeeper,      // Assign bankKeeper
		accountKeeper: accountKeeper, // Assign accountKeeper
		TokenSeq: collections.NewSequence(sb, types.TokenKeyPrefix+"_seq"),
		Token: collections.NewMap(sb, types.TokenKeyPrefix, collections.Uint64Key, sdk.NewProtoValue(types.Token{})),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) Logger(ctx context.Context) log.Logger { // Use context.Context for collections
	sdkCtx := sdk.UnwrapSDKContext(ctx) // Unwrap to sdk.Context for logger if needed
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetTokenBySymbol: Iterate through tokens to find by symbol.
// In a production environment, for better performance, you would add a secondary
// index (e.g., collections.Map[string, uint64]) to map symbols to token IDs.
func (k Keeper) GetTokenBySymbol(ctx context.Context, symbol string) (val types.Token, found bool) {
	// Iterate through all tokens. This is inefficient for a large number of tokens.
	// For simplicity in this tutorial, we'll use this method.
	// For a scalable solution, consider a collections.Map[string, uint64] to map symbols to IDs.
	err := k.Token.Walk(ctx, nil, func(key uint64, token types.Token) (stop bool, err error) {
		if token.Symbol == symbol {
			val = token
			found = true
			return true, nil // Stop iteration
		}
		return false, nil // Continue iteration
	})
	if err != nil {
		// Log error or handle appropriately, but don't return error in signature
		// if the goal is just to return found status.
		k.Logger(ctx).Error("error walking tokens to find by symbol", "error", err)
	}
	return val, found
}
