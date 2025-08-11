package token

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"omnis/x/token/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListToken",
					Use:       "list-token",
					Short:     "List all token",
				},
				{
					RpcMethod:      "GetToken",
					Use:            "get-token [id]",
					Short:          "Gets a token by id",
					Alias:          []string{"show-token"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateToken",
					Use:            "create-token [name] [symbol] [decimals] [total-supply] [metadata]",
					Short:          "Create token",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "name"}, {ProtoField: "symbol"}, {ProtoField: "decimals"}, {ProtoField: "total_supply"}, {ProtoField: "metadata"}},
				},
				{
					RpcMethod:      "UpdateToken",
					Use:            "update-token [id] [name] [symbol] [decimals] [total-supply] [metadata]",
					Short:          "Update token",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "name"}, {ProtoField: "symbol"}, {ProtoField: "decimals"}, {ProtoField: "total_supply"}, {ProtoField: "metadata"}},
				},
				{
					RpcMethod:      "DeleteToken",
					Use:            "delete-token [id]",
					Short:          "Delete token",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
