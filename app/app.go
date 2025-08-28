package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"cosmossdk.io/server/api"
	serverconfig "cosmossdk.io/server/config"
	servertypes "cosmossdk.io/server/types"
	"cosmossdk.io/store/types"
	"cosmossdk.io/x/auth"
	authkeeper "cosmossdk.io/x/auth/keeper"
	authtypes "cosmossdk.io/x/auth/types"
	"cosmossdk.io/x/auth/tx/config"
	authz "cosmossdk.io/x/authz"
	authzkeeper "cosmossdk.io/x/authz/keeper"
	authzmodule "cosmossdk.io/x/authz/module"
	"cosmossdk.io/x/bank"
	bankkeeper "cosmossdk.io/x/bank/keeper"
	banktypes "cosmossdk.io/x/bank/types"
	"cosmossdk.io/x/consensus"
	consensusparamkeeper "cosmossdk.io/x/consensus/keeper"
	"cosmossdk.io/x/crisis"
	crisiskeeper "cosmossdk.io/x/crisis/keeper"
	crisistypes "cosmossdk.io/x/crisis/types"
	"cosmossdk.io/x/distribution"
	distrkeeper "cosmossdk.io/x/distribution/keeper"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/genutil"
	genutiltypes "cosmossdk.io/x/genutil/types"
	"cosmossdk.io/x/gov"
	govkeeper "cosmossdk.io/x/gov/keeper"
	govtypes "cosmossdk.io/x/gov/types"
	"cosmossdk.io/x/group"
	groupkeeper "cosmossdk.io/x/group/keeper"
	groupmodule "cosmossdk.io/x/group/module"
	"cosmossdk.io/x/mint"
	mintkeeper "cosmossdk.io/x/mint/keeper"
	minttypes "cosmossdk.io/x/mint/types"
	"cosmossdk.io/x/slashing"
	slashingkeeper "cosmossdk.io/x/slashing/keeper"
	"cosmossdk.io/x/staking"
	stakingkeeper "cosmossdk.io/x/staking/keeper"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/legacy"
	"github.com/cosmos/cosmos-sdk/server/rosetta"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/omnis-org/omnis/x/omnis"
	omniskeeper "github.com/omnis-org/omnis/x/omnis/keeper"
	omnistypes "github.com/omnis-org/omnis/x/omnis/types"
	"github.com/omnis-org/omnis/x/token"
	tokenkeeper "github.com/omnis-org/omnis/x/token/keeper"
	tokentypes "github.com/omnis-org/omnis/x/token/types"

	// IBC
	ibcclient "github.com/cosmos/ibc-go/v10/modules/core/02-client"
	ibcclientkeeper "github.com/cosmos/ibc-go/v10/modules/core/02-client/keeper"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	capabilitykeeper "github.com/cosmos/ibc-go/v10/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/v10/modules/capability/types"
	ibcporttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/01-tendermint"
)

var (
	_ servertypes.Application = (*OmnisApp)(nil)

	// DefaultNodeHome is the default home directory for the chain's node daemon.
	DefaultNodeHome string
)

// App extends the ABCI application for the Cosmos SDK.
type OmnisApp struct {
	*runtime.App

	// LegacyBaseApp is the ABCI application that is wrapped by the new runtime.
	LegacyBaseApp *legacy.BaseApp

	// Codec
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	cdc               codec.Codec
	interfaceRegistry types.InterfaceRegistry

	// Keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistributionKeeper    distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	FeegrantKeeper        feegrantkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// IBC
	IBCKeeper *ibckeeper.Keeper
	// Scoped keepers
	ScopedIBCKeeper       capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper  capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper

	// Custom modules
	OmnisKeeper omniskeeper.Keeper
	TokenKeeper tokenkeeper.Keeper

	// Module Manager
	ModuleManager *appmodule.Manager

	// Home path
	HomePath string

	// Config
	Viper *viper.Viper

	// Ticker for the application
	Ticker *time.Ticker

	// Logger
	Logger log.Logger
}

// NewOmnisApp initializes the application.
func NewOmnisApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
) *OmnisApp {
	// A new Go module, `depinject`, is used to handle dependency injection.
	// This helps with managing the keepers and their dependencies.
	var (
		app *runtime.App
		// Autocli is a new way to generate CLI commands for your modules.
		autocliOpts autocli.AppOptions
		// Module keepers
		accountKeeper         authkeeper.AccountKeeper
		bankKeeper            bankkeeper.Keeper
		stakingKeeper         *stakingkeeper.Keeper
		slashingKeeper        slashingkeeper.Keeper
		mintKeeper            mintkeeper.Keeper
		distributionKeeper    distrkeeper.Keeper
		govKeeper             govkeeper.Keeper
		crisisKeeper          *crisiskeeper.Keeper
		upgradeKeeper         *upgradekeeper.Keeper
		authzKeeper           authzkeeper.AuthzKeeper
		evidenceKeeper        evidencekeeper.Keeper
		feegrantKeeper        feegrantkeeper.Keeper
		groupKeeper           groupkeeper.Keeper
		consensusParamsKeeper consensusparamkeeper.Keeper

		// Custom keepers
		omnisKeeper omniskeeper.Keeper
		tokenKeeper tokenkeeper.Keeper

		// IBC keepers
		ibcKeeper *ibckeeper.Keeper

		// New runtime configuration is done here.
		// A new method, `app.Build`, is used to construct the application.
	)

	// depinject is used to configure the modules and their dependencies.
	// It replaces the old, manual way of wiring up keepers in the constructor.
	// A lot of the `app.go` logic is now handled by this.
	depinject.Make(
		depinject.Configs(
			serverconfig.AppConfig,
			depinject.Provide(
				// Define all custom modules here.
				// This is where you would provide your x/token and x/omnis keepers.
				func(
					m map[string]appmodule.AppModule,
					cfg depinject.Config,
					logger log.Logger,
					cdc codec.Codec,
					homePath string,
					ac authtypes.AccountKeeper,
					bk banktypes.BankKeeper,
				) (map[string]appmodule.AppModule, omniskeeper.Keeper, tokenkeeper.Keeper) {
					omnisKeeper = omniskeeper.NewKeeper(
						cdc,
						logger,
						homePath,
					)

					tokenKeeper = tokenkeeper.NewKeeper(
						cdc,
						logger,
						homePath,
					)

					m["omnis"] = omnis.NewAppModule(omnisKeeper)
					m["token"] = token.NewAppModule(tokenKeeper)

					return m, omnisKeeper, tokenKeeper
				},
			),
		),
		&app,
		&autocliOpts,
		&accountKeeper,
		&bankKeeper,
		&stakingKeeper,
		&slashingKeeper,
		&mintKeeper,
		&distributionKeeper,
		&govKeeper,
		&crisisKeeper,
		&upgradeKeeper,
		&authzKeeper,
		&evidenceKeeper,
		&feegrantKeeper,
		&groupKeeper,
		&consensusParamsKeeper,
		&omnisKeeper,
		&tokenKeeper,
		&ibcKeeper,
	)

	// The app struct is now initialized with the new runtime.
	legacyBaseApp := legacy.NewBaseApp(
		"omnis",
		logger,
		db,
		app.TxConfig.TxDecoder(),
		nil, // Replaced by runtime.App
		nil, // Replaced by runtime.App
	)
	legacyBaseApp.SetCommitMultiStore(app.CommitMultiStore())
	legacyBaseApp.SetRouter(app.Router())
	legacyBaseApp.SetQueryRouter(app.QueryRouter())

	// Init the IBC Keeper
	ibcKeeper = ibckeeper.NewKeeper(
		app.AppCodec(),
		app.State,
		app.ScopedIBCKeeper,
		app.ScopedTransferKeeper,
		nil,
	)

	// The IBC client keeper is added here.
	ibcClientKeeper := ibcclientkeeper.NewKeeper(
		app.AppCodec(),
		app.State,
		app.ScopedIBCKeeper,
		ibcKeeper.ConnectionKeeper,
		ibcKeeper.ChannelKeeper,
	)

	app.ModuleManager.RegisterModules(
		auth.NewAppModule(accountKeeper, nil),
		bank.NewAppModule(bankKeeper, nil),
		staking.NewAppModule(stakingKeeper, accountKeeper, bankKeeper, nil),
		mint.NewAppModule(mintKeeper, accountKeeper, bankKeeper, nil),
		distribution.NewAppModule(distributionKeeper, accountKeeper, bankKeeper, nil),
		gov.NewAppModule(govKeeper, accountKeeper, bankKeeper, nil),
		crisis.NewAppModule(crisisKeeper, nil),
		slashing.NewAppModule(slashingkeeper, accountKeeper, bankKeeper, nil),
		feegrantmodule.NewAppModule(accountKeeper, bankKeeper, feegrantkeeper, nil),
		upgrade.NewAppModule(upgradekeeper),
		evidence.NewAppModule(evidencekeeper),
		authzmodule.NewAppModule(authzkeeper, accountKeeper, nil),
		groupmodule.NewAppModule(groupkeeper, nil),
		consensus.NewAppModule(consensusParamsKeeper),

		// IBC modules
		ibcclient.NewAppModule(ibcClientKeeper),
		ibctm.NewAppModule(),
	)

	// Create new OmnisApp struct with all keepers
	newApp := &OmnisApp{
		App:                   app,
		LegacyBaseApp:         legacyBaseApp,
		appCodec:              app.AppCodec(),
		interfaceRegistry:     app.InterfaceRegistry(),
		AccountKeeper:         accountKeeper,
		BankKeeper:            bankKeeper,
		StakingKeeper:         stakingKeeper,
		SlashingKeeper:        slashingKeeper,
		MintKeeper:            mintKeeper,
		DistributionKeeper:    distributionKeeper,
		GovKeeper:             govKeeper,
		CrisisKeeper:          crisisKeeper,
		UpgradeKeeper:         upgradeKeeper,
		AuthzKeeper:           authzkeeper,
		EvidenceKeeper:        evidenceKeeper,
		FeegrantKeeper:        feegrantKeeper,
		GroupKeeper:           groupKeeper,
		ConsensusParamsKeeper: consensusParamsKeeper,
		IBCKeeper:             ibcKeeper,
		OmnisKeeper:           omnisKeeper,
		TokenKeeper:           tokenKeeper,
		HomePath:              appOpts.Get(flags.FlagHome).(string),
		Logger:                logger,
	}

	// Set handlers for the app.
	newApp.MountStores()

	// The app is now ready.
	return newApp
}

// MountStores mounts all the store keys.
func (app *OmnisApp) MountStores() {
	// Your custom store keys and other app setup logic can go here.
	// This is where you would set up your key-value stores for the modules.
	keys := make([]types.StoreKey, 0)
	for _, module := range app.ModuleManager.Modules {
		keys = append(keys, module.StoreKeys()...)
	}

	// Mount your custom module store keys
	keys = append(keys, omnistypes.StoreKey, tokentypes.StoreKey)

	for _, key := range keys {
		app.LegacyBaseApp.MountStore(key, types.StoreKeyType)
	}
}

// Name returns the name of the App.
func (app *OmnisApp) Name() string { return app.App.Name() }

// GetBaseApp is for testing purposes.
func (app *OmnisApp) GetBaseApp() *legacy.BaseApp { return app.LegacyBaseApp }

// PreBlocker returns the PreBlocker.
func (app *OmnisApp) PreBlocker(ctx context.Context) error {
	return app.App.PreBlocker(ctx)
}

// BeginBlocker returns the BeginBlocker.
func (app *OmnisApp) BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) (abci.ResponseBeginBlock, error) {
	return app.App.BeginBlocker(ctx, req)
}

// EndBlocker returns the EndBlocker.
func (app *OmnisApp) EndBlocker(ctx context.Context, req abci.RequestEndBlock) (abci.ResponseEndBlock, error) {
	return app.App.EndBlocker(ctx, req)
}

// DeliverTx returns the DeliverTx.
func (app *OmnisApp) DeliverTx(req abci.RequestDeliverTx) (abci.ResponseDeliverTx, error) {
	return app.LegacyBaseApp.DeliverTx(req)
}

// LoadHeight returns the LoadHeight.
func (app *OmnisApp) LoadHeight(height int64) error {
	return app.App.LoadHeight(height)
}

// LastCommitId returns the last commit ID.
func (app *OmnisApp) LastCommitId() types.CommitID {
	return app.LegacyBaseApp.LastCommitID()
}

// LastBlockHeight returns the last block height.
func (app *OmnisApp) LastBlockHeight() int64 {
	return app.LegacyBaseApp.LastBlockHeight()
}

// CheckTx returns the CheckTx.
func (app *OmnisApp) CheckTx(req abci.RequestCheckTx) (abci.ResponseCheckTx, error) {
	return app.LegacyBaseApp.CheckTx(req)
}
