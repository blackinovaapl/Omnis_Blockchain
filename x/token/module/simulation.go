package token

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"omnis/testutil/sample"
	tokensimulation "omnis/x/token/simulation"
	"omnis/x/token/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenGenesis := types.GenesisState{
		Params:    types.DefaultParams(),
		TokenList: []types.Token{{Id: 0, Creator: sample.AccAddress()}, {Id: 1, Creator: sample.AccAddress()}}, TokenCount: 2,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateToken          = "op_weight_msg_token"
		defaultWeightMsgCreateToken int = 100
	)

	var weightMsgCreateToken int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateToken, &weightMsgCreateToken, nil,
		func(_ *rand.Rand) {
			weightMsgCreateToken = defaultWeightMsgCreateToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateToken,
		tokensimulation.SimulateMsgCreateToken(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateToken          = "op_weight_msg_token"
		defaultWeightMsgUpdateToken int = 100
	)

	var weightMsgUpdateToken int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateToken, &weightMsgUpdateToken, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateToken = defaultWeightMsgUpdateToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateToken,
		tokensimulation.SimulateMsgUpdateToken(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteToken          = "op_weight_msg_token"
		defaultWeightMsgDeleteToken int = 100
	)

	var weightMsgDeleteToken int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteToken, &weightMsgDeleteToken, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteToken = defaultWeightMsgDeleteToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteToken,
		tokensimulation.SimulateMsgDeleteToken(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
