package keeper_test

import (
	"testing"

	"omnis/x/token/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:     types.DefaultParams(),
		TokenList:  []types.Token{{Id: 0}, {Id: 1}},
		TokenCount: 2,
	}
	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.TokenList, got.TokenList)
	require.Equal(t, genesisState.TokenCount, got.TokenCount)

}
