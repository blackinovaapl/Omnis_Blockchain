package keeper

import (
	"context"

	"omnis/x/token/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.TokenList {
		if err := k.Token.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.TokenSeq.Set(ctx, genState.TokenCount); err != nil {
		return err
	}
	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = k.Token.Walk(ctx, nil, func(key uint64, elem types.Token) (bool, error) {
		genesis.TokenList = append(genesis.TokenList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.TokenCount, err = k.TokenSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
