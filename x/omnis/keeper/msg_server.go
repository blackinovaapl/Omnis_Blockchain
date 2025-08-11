package keeper

import (
	"context"
	"fmt" // For error formatting

	sdk "github.com/cosmos/cosmos-sdk/types" // For sdk.Context and sdk.AccAddress
	"omnis/x/omnis/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateItem implements types.MsgServer
func (k msgServer) CreateItem(ctx context.Context, msg *types.MsgCreateItem) (*types.MsgCreateItemResponse, error) {
    sdkCtx := sdk.UnwrapSDKContext(ctx)

    creatorAddr, err := k.addressCodec.StringToBytes(msg.Creator)
    if err != nil {
        return nil, err
    }
    // Basic validation (e.g., ensure creator is valid, item doesn't already exist)
    // You'll likely need to generate a unique ID for the item.
    // For simplicity, let's assume item.Id is passed or generated here.
    // In a real application, you'd have a counter or UUID generator.
    if msg.Id == "" { // Or generate a new ID if it's not provided by the user
        return nil, fmt.Errorf("item ID cannot be empty")
    }

    // Check if item already exists
    _, err = k.Keeper.GetItem(ctx, msg.Id)
    if err == nil { // No error means it exists
        return nil, fmt.Errorf("item with ID %s already exists", msg.Id)
    }
    // If it's collections.ErrNotFound, then it's good to proceed.
    // Otherwise, it's an unexpected error.
    if !errors.Is(err, collections.ErrNotFound) {
        return nil, err
    }


    item := types.Item{
        Id:    msg.Id,
        Name:  msg.Name,
        Owner: msg.Creator, // Assuming owner is the creator
        // ... other fields from your types.Item
    }

    if err := k.Keeper.SetItem(ctx, item); err != nil {
        return nil, err
    }

    // Emit events (important for clients and indexers)
    sdkCtx.EventManager().EmitEvents(sdk.Events{
        sdk.NewEvent(
            types.EventTypeCreateItem,
            sdk.NewAttribute(types.AttributeKeyItemID, item.Id),
            sdk.NewAttribute(types.AttributeKeyItemName, item.Name),
            sdk.NewAttribute(types.AttributeKeyCreator, item.Owner),
        ),
    })

    return &types.MsgCreateItemResponse{}, nil
}

// UpdateItem implements types.MsgServer
func (k msgServer) UpdateItem(ctx context.Context, msg *types.MsgUpdateItem) (*types.MsgUpdateItemResponse, error) {
    sdkCtx := sdk.UnwrapSDKContext(ctx)

    // Get the existing item
    item, err := k.Keeper.GetItem(ctx, msg.Id)
    if err != nil {
        return nil, fmt.Errorf("item with ID %s not found: %w", msg.Id, err)
    }

    // Check if the sender is the owner of the item
    if item.Owner != msg.Creator {
        return nil, fmt.Errorf("unauthorized: sender is not the item owner")
    }

    // Update fields
    item.Name = msg.NewName // Assuming NewName is a field in MsgUpdateItem
    // ... update other fields as per your MsgUpdateItem definition

    if err := k.Keeper.SetItem(ctx, item); err != nil {
        return nil, err
    }

    sdkCtx.EventManager().EmitEvents(sdk.Events{
        sdk.NewEvent(
            types.EventTypeUpdateItem,
            sdk.NewAttribute(types.AttributeKeyItemID, item.Id),
            sdk.NewAttribute(types.AttributeKeyNewItemName, item.Name),
            sdk.NewAttribute(types.AttributeKeyCreator, item.Owner),
        ),
    })

    return &types.MsgUpdateItemResponse{}, nil
}

// DeleteItem implements types.MsgServer
func (k msgServer) DeleteItem(ctx context.Context, msg *types.MsgDeleteItem) (*types.MsgDeleteItemResponse, error) {
    sdkCtx := sdk.UnwrapSDKContext(ctx)

    item, err := k.Keeper.GetItem(ctx, msg.Id)
    if err != nil {
        return nil, fmt.Errorf("item with ID %s not found: %w", msg.Id, err)
    }

    if item.Owner != msg.Creator {
        return nil, fmt.Errorf("unauthorized: sender is not the item owner")
    }

    if err := k.Keeper.DeleteItem(ctx, msg.Id); err != nil {
        return nil, err
    }

    sdkCtx.EventManager().EmitEvents(sdk.Events{
        sdk.NewEvent(
            types.EventTypeDeleteItem,
            sdk.NewAttribute(types.AttributeKeyItemID, item.Id),
            sdk.NewAttribute(types.AttributeKeyCreator, item.Owner),
        ),
    })

    return &types.MsgDeleteItemResponse{}, nil
}