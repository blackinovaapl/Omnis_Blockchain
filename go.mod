module github.com/omnis-org/omnis

go 1.24

require (
	github.com/cometbft/cometbft v0.37.5
	github.com/cosmos/cosmos-sdk v0.50.0
	github.com/cosmos/cosmos-db v1.0.2
	github.com/cosmos/ibc-go/v10 v10.2.0
	github.com/cosmos/gogoproto v1.4.10
	github.com/gorilla/mux v1.8.1
	github.com/omnis-org/omnis/x/omnis v0.0.0
	github.com/omnis-org/omnis/x/token v0.0.0
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.6.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.16.0
	golang.org/x/exp v0.0.0-20230510235704-dd9e11501b44
)

// This replace directive is crucial. It tells Go to use the
// 'cosmossdk.io/cosmos-sdk' module for all `github.com/cosmos/cosmos-sdk`
// imports. The previous path was incorrect. This path will be found
// in the Go module proxy and resolve the 404 error.
replace cosmossdk.io/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.50.0

// These replace directives are for local development of your modules.
replace github.com/omnis-org/omnis/x/omnis => ./x/omnis
replace github.com/omnis-org/omnis/x/token => ./x/token
