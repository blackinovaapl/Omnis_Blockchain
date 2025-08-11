// x/omnis/keeper/query_server.go (new file or add to query_params.go and rename)
package keeper

import (
	"context"
	"errors" // For collections.ErrNotFound
	"fmt"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"omnis/x/omnis/types"
)

type queryServer struct {
	Keeper // The query server needs access to the Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{} // This confirms it implements the interface

// GetItem handles the QueryGetItemRequest
func (q queryServer) GetItem(ctx context.Context, req *types.QueryGetItemRequest) (*types.QueryGetItemResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    item, err := q.Keeper.GetItem(ctx, req.Id)
    if err != nil {
        if errors.Is(err, collections.ErrNotFound) {
            return nil, status.Error(codes.NotFound, fmt.Sprintf("item with ID %s not found", req.Id))
        }
        return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get item: %s", err.Error()))
    }

    return &types.QueryGetItemResponse{Item: item}, nil
}

// AllItems handles the QueryAllItemsRequest
func (q queryServer) AllItems(ctx context.Context, req *types.QueryAllItemsRequest) (*types.QueryAllItemsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    var items []types.Item
    // Iterate over all items using the helper method from the keeper
    err := q.Keeper.IterateItems(ctx, func(id string, item types.Item) (bool, error) {
        items = append(items, item)
        return false, nil // Continue iteration
    })
    if err != nil {
        return nil, status.Error(codes.Internal, fmt.Sprintf("failed to iterate items: %s", err.Error()))
    }

    return &types.QueryAllItemsResponse{Items: items}, nil
}