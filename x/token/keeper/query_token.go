package keeper

import (
	"context"
	"errors"

	"omnis/x/token/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListToken(ctx context.Context, req *types.QueryAllTokenRequest) (*types.QueryAllTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokens, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Token,
		req.Pagination,
		func(_ uint64, value types.Token) (types.Token, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTokenResponse{Token: tokens, Pagination: pageRes}, nil
}

func (q queryServer) GetToken(ctx context.Context, req *types.QueryGetTokenRequest) (*types.QueryGetTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	token, err := q.k.Token.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetTokenResponse{Token: token}, nil
}
