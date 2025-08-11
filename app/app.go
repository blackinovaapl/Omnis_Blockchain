package app

import (
	"io"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	circuit "cosmossdk.io/x/circuit"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	upgrade "cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Cosmos SDK modules using cosmossdk.io paths
	auth "cosmossdk.io/x/auth"
	authkeeper "cosmossdk.io/x/auth/keeper"
	authsims "cosmossdk.io/x/auth/simulation"
	authtypes "cosmossdk.io/x/auth/types"
	authz "cosmossdk.io/x/authz"
	authzkeeper "cosmossdk.io/x/authz/keeper"
	bank "cosmossdk.io/x/bank"
	bankkeeper "cosmossdk.io/x/bank/keeper"
	banktypes "cosmossdk.io/x/bank/types"
	consensus "cosmossdk.io/x/consensus"
	consensuskeeper "cosmossdk.io/x/consensus/keeper"
	consensustypes "cosmossdk.io/x/consensus/types"
	distr "cosmossdk.io/x/distribution"
	distrkeeper "cosmossdk.io/x/distribution/keeper"
	distrtypes "cosmossdk.io/x/distribution/types"
	evidence "cosmossdk.io/x/evidence" // Corrected import
	evidencetypes "cosmossdk.io/x/evidence/types" // Corrected import
	feegrant "cosmossdk.io/x/feegrant" // Corrected import
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper" // Added
	genutil "github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	gov "cosmossdk.io/x/gov"
	govkeeper "cosmossdk.io/x/gov/keeper"
	govtypes "cosmossdk.io/x/gov/types"
	mint "cosmossdk.io/x/mint"
	mintkeeper "cosmossdk.io/x/mint/keeper"
	minttypes "cosmossdk.io/x/mint/types"
	params "cosmossdk.io/x/params"
	paramskeeper "cosmossdk.io/x/params/keeper"
	paramstypes "cosmossdk.io/x/params/types"
	slashing "cosmossdk.io/x/slashing"
	slashingkeeper "cosmossdk.io/x/slashing/keeper"
	slashingtypes "cosmossdk.io/x/slashing/types"
	staking "cosmossdk.io/x/staking"
	stakingkeeper "cosmossdk.io/x/staking/keeper"
	stakingtypes "cosmossdk.io/x/staking/types"


	// IBC modules using github.com/cosmos/ibc-go/v10 paths
	
	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibchost "github.com/cosmos/ibc-go/v10/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	ica "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/keeper"
	icahostkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/keeper"
	ibctransfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"

	routetypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types" // For port keeper route

	"omnis/docs"
	omnismodule "omnis/x/omnis"
	omnismodulekeeper "omnis/x/omnis/keeper"

	"omnis/x/token"
	tokentypes "omnis/x/token/types"
	tokenkeeper "omnis/x/token/keeper"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1" // Keep this as it is
)

const (
	// Name is the name of the application.
	Name = "omnis"
	// AccountAddressPrefix is the prefix for accounts addresses.
	AccountAddressPrefix = "cosmos"
	// ChainCoinType is the coin type of the chain.
	ChainCoinType = 118
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
	// App extends an ABCI application, but with most of its parameters exported.
	// They are exported for convenience in creating helper functions, as object
	// capabilities aren't needed for testing.
	_ runtime.App              = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

type App struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keepers
	// only keepers required by the app are exposed
	// the list of all modules is available in the app_config
	AuthKeeper          authkeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	StakingKeeper       *stakingkeeper.Keeper
	SlashingKeeper      slashingkeeper.Keeper
	MintKeeper          mintkeeper.Keeper
	DistrKeeper         distrkeeper.Keeper
	GovKeeper           *govkeeper.Keeper
	UpgradeKeeper       *upgradekeeper.Keeper
	AuthzKeeper         authzkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper
	CircuitBreakerKeeper circuitkeeper.Keeper
	ParamsKeeper        paramskeeper.Keeper
	FeeGrantKeeper      feegrantkeeper.Keeper // Added

	// ibc keepers
	IBCKeeper           *ibckeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper

	OmnisKeeper omnismodulekeeper.Keeper
	TokenKeeper tokenkeeper.Keeper
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	sdk.DefaultBondDenom = "stake"
	var err error
	clienthelpers.EnvPrefix = Name
	DefaultNodeHome, err = clienthelpers.GetNodeDirectory(Name)
	if err != nil {
		panic(err)
	}
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appConfig,
		depinject.Supply(
			// supply custom module basics
			map[string]appmodule.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			},
		),
	)
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appopts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	app := &App{}
	appBuilder := &runtime.AppBuilder{}

	// merge the AppConfig and other configuration in one config
	appConfig := depinject.Configs(
		AppConfig(),
		depinject.Supply(
			// supply app options
			appopts,
			logger,
			// supply logger
			// here alternative options can be supplied to the DI container.
			// those options can be used f. e to override the default behavior of some modules.
			// for instance supplying a custom address codec for not using bech32 addresses.
			// read the depinject documentation and depinject module wiring for more information
			// on available options and how to use them.
		),
	)

	var appmodules map[string]appmodule.AppModule
	if err := depinject.Inject(appConfig,
		&appBuilder,
		&appmodules,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AuthKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.MintKeeper,
		&app.DistrKeeper,
		&app.GovKeeper,
		&app.UpgradeKeeper,
		&app.AuthzKeeper,
		&app.ConsensusParamsKeeper,
		&app.CircuitBreakerKeeper,
		&app.ParamsKeeper,
		&app.OmnisKeeper,
		&app.TokenKeeper,
		&app.FeeGrantKeeper, // <--- Added FeeGrantKeeper to inject
		// this line is used by starport scaffolding # stargate/app/keeperDeclaration
	); err != nil {
		panic(err)
	}

	// Initialize TokenKeeper after other core keepers are available
	app.TokenKeeper = tokenkeeper.NewKeeper(
		app.AuthKeeper.AddressCodec(),
		app.BaseApp.StoreService(),
		app.BankKeeper,
		app.AuthKeeper,
	)

	// add to default baseapp options
	// enable optimistic execution
	baseAppOptions = append(baseAppOptions, baseapp.SetOptimisticExecution())

	// build app
	app.App = appBuilder.Build(db, traceStore,
		baseAppOptions...,
	)

	// Define the module manager with all your modules
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.appCodec,
			app.GetEnvProposerAPIRouter(),
			app.txConfig,
			app.AuthKeeper,
			app.BankKeeper,
			app.ConsensusParamsKeeper,
		),


		   // create the simulation manager and define the order of the modules for deterministic simulations
    overridemodules:= map[string]module.AppModuleSimulation{
        authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AuthKeeper, authsims.RandomGenesisAccounts, app.GetTxConfig()),
    }
    app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overridemodules)
    app.sm.RegisterStoreDecoders()


	
		auth.NewAppModule(app.appCodec, app.AuthKeeper, authsims.RandomGenesisAccounts, app.GetTxConfig()),
		bank.NewAppModule(app.appCodec, app.BankKeeper, app.AuthKeeper.AccountKeeper()),
		staking.NewAppModule(app.appCodec, app.StakingKeeper, app.AuthKeeper.AccountKeeper(), app.BankKeeper),
		mint.NewAppModule(app.appCodec, app.MintKeeper, app.AuthKeeper.AccountKeeper(), nil, app.BankKeeper),
		distr.NewAppModule(app.appCodec, app.DistrKeeper, app.AuthKeeper.AccountKeeper(), app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(app.appCodec, app.SlashingKeeper, app.AuthKeeper.AccountKeeper(), app.BankKeeper, app.StakingKeeper),
		gov.NewAppModule(app.appCodec, app.GovKeeper, app.AuthKeeper.AccountKeeper(), app.BankKeeper, app.GetModuleManager(), app.txConfig),
		params.NewAppModule(app.ParamsKeeper),
		circuit.NewAppModule(app.appCodec, app.CircuitBreakerKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		authz.NewAppModule(app.appCodec, app.AuthzKeeper, app.AuthKeeper.AccountKeeper(), app.BankKeeper, app.interfaceRegistry),
		feegrant.NewAppModule(app.appCodec, app.BankKeeper, app.AuthKeeper.AccountKeeper(), app.interfaceRegistry),
		evidence.NewAppModule(app.AuthKeeper.AccountKeeper(), app.BankKeeper, app.SlashingKeeper), // Added Evidence Module
		consensus.NewAppModule(app.appCodec, app.ConsensusParamsKeeper), // Added Consensus Module
		// IBC modules
		ibctransfer.NewAppModule(app.TransferKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		omnismodule.NewAppModule(app.appCodec, app.OmnisKeeper, app.AuthKeeper.AccountKeeper(), app.BankKeeper),
		token.NewAppModule(app.appCodec, app.TokenKeeper, app.AuthKeeper.AccountKeeper(), app.BankKeeper),
		// this line is used by starport scaffolding # stargate/app/module
	)

	// create the simulation manager and define the order of the modules for deterministic simulations
	overridemodules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AuthKeeper, authsims.RandomGenesisAccounts, app.GetTxConfig()),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overridemodules)
	app.sm.RegisterStoreDecoders()

	// A custom InitChainer sets if extra pre-init-genesis logic is required.
	// This is necessary for manually registered modules that do not support app wiring.
	// Manually set the module version map as shown below.
	// The upgrade module will automatically handle de-duplication of the module version map.

	// Corrected SetInitChainer:
	app.SetInitChainer(func(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
		if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil { // Corrected: app.mm.GetVersionMap()
			return nil, err
		}
		return app.App.InitChainer(ctx, req)
	})


	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

// GetSubspace returns a param subspace for a given module name.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// LegacyAmino returns App's amino codec.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns App's app codec.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns App's InterfaceRegistry.
func (app *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns App's TxConfig
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}

	// register app's OpenAPI routes.
	docs.RegisterOpenAPIService(Name, apiSvr.Router)
}

// GetMaccPerms returns a copy of the module account permissions
// NOTE: This is solely to be used for testing purposes.
func GetMaccPerms() map[string][]string {
	dup := make(map[string][]string)
	for acc, perms := range maccPerms {
		dup[acc] = perms
	}
	return dup
}

// maccPerms are the permissions for module accounts.
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	tokentypes.ModuleName:          {authtypes.Minter, authtypes.Burner},
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	result := make(map[string]bool)
	for acc := range GetMaccPerms() {
		result[authtypes.NewModuleAddress(acc).String()] = true
	}
	return result
}

// For Store and Transient Keys, these are typically defined in the `New` function.
// Let's ensure they are correct inside `func New(...)`

/*
	// ... inside New func, after depinject.Inject and before appBuilder.Build()

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramtypes.StoreKey,
		upgradetypes.StoreKey,
		feegrant.StoreKey,
		evidencetypes.StoreKey,
		ibchost.StoreKey,
		ibctransfertypes.StoreKey,

		routetypes.StoreKey, // This might be `porttypes.StoreKey` depending on your IBC version
		consensustypes.StoreKey,
		tokentypes.StoreKey,
		autocliv1.StoreKey,


	)
	memKeys := sdk.NewTransientStoreKeys(
		paramtypes.TStoreKey,
	
		tokentypes.MemStoreKey,
	)

	// ... continue with appBuilder.Build() etc.

*/
