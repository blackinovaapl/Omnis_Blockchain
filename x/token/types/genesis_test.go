package types_test

import (
	"testing"

	"omnis/x/token/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &types.GenesisState{TokenList: []types.Token{{Id: 0}, {Id: 1}}, TokenCount: 2}, valid: true,
		}, {
			desc: "duplicated token",
			genState: &types.GenesisState{
				TokenList: []types.Token{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		}, {
			desc: "invalid token count",
			genState: &types.GenesisState{
				TokenList: []types.Token{
					{
						Id: 1,
					},
				},
				TokenCount: 0,
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
