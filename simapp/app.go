package simapp

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/depinject"

	cmtDb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side effects
	// Authority
	"github.com/noble-assets/paramauthority/x/authority"
	authorityKeeper "github.com/noble-assets/paramauthority/x/authority/keeper"
	// Bank
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Capability
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	// Consensus
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	// GenUtil
	genUtil "github.com/cosmos/cosmos-sdk/x/genutil"
	genUtilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	// IBC Client
	ibcClientSolomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibcClientTendermint "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	// IBC Core
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcKeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	// IBC Fee
	ibcFee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcFeeKeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	// IBC Transfer
	ibcTransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	// ICA
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	// ICA Controller
	icaControllerKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	// ICA Host
	icaHostKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	// Params
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	// Staking
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	// Upgrade
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
)

var (
	DefaultNodeHome string

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		consensus.AppModuleBasic{},
		genUtil.NewAppModuleBasic(genUtilTypes.DefaultMessageValidator),
		params.AppModuleBasic{},
		staking.AppModuleBasic{},
		upgrade.AppModuleBasic{},

		ibc.AppModuleBasic{},
		ibcClientSolomachine.AppModuleBasic{},
		ibcClientTendermint.AppModuleBasic{},
		ibcFee.AppModuleBasic{},
		ibcTransfer.AppModuleBasic{},
		ica.AppModuleBasic{},

		authority.AppModuleBasic{},
	)
)

var (
	_ runtime.AppI            = (*SimApp)(nil)
	_ serverTypes.Application = (*SimApp)(nil)
)

type SimApp struct {
	*runtime.App

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codecTypes.InterfaceRegistry

	// Cosmos SDK Keepers
	AccountKeeper    authKeeper.AccountKeeper
	BankKeeper       bankKeeper.Keeper
	CapabilityKeeper *capabilityKeeper.Keeper
	ConsensusKeeper  consensusKeeper.Keeper
	ParamsKeeper     paramsKeeper.Keeper
	StakingKeeper    *stakingKeeper.Keeper
	UpgradeKeeper    *upgradeKeeper.Keeper

	// IBC Keepers
	IBCKeeper           *ibcKeeper.Keeper
	IBCFeeKeeper        ibcFeeKeeper.Keeper
	IBCTransferKeeper   ibcTransferKeeper.Keeper
	ICAControllerKeeper icaControllerKeeper.Keeper
	ICAHostKeeper       icaHostKeeper.Keeper

	// Custom Keepers
	AuthorityKeeper *authorityKeeper.Keeper

	// Scoped Keepers (for IBC)
	ScopedIBCKeeper           capabilityKeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilityKeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilityKeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilityKeeper.ScopedKeeper
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".simapp")
}

func NewSimApp(
	logger log.Logger,
	db cmtDb.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts serverTypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *SimApp {
	var (
		app        = &SimApp{}
		appBuilder *runtime.AppBuilder

		appConfig = depinject.Configs(
			AppConfig,
			depinject.Supply(appOpts),
		)
	)

	if err := depinject.Inject(appConfig,
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.CapabilityKeeper,
		&app.ConsensusKeeper,
		&app.ParamsKeeper,
		&app.StakingKeeper,
		&app.UpgradeKeeper,
		&app.AuthorityKeeper,
	); err != nil {
		panic(err)
	}

	app.App = appBuilder.Build(logger, db, traceStore, baseAppOptions...)

	// Registers all modules that don't use App Wiring (e.g. IBC).
	app.RegisterLegacyModules()
	// Registers all proposals handlers that are using v1beta1 governance.
	app.RegisterLegacyRouter()

	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

func (app *SimApp) GetSubspace(moduleName string) paramsTypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// ------------------------------- runtime.AppI --------------------------------

func (app *SimApp) AppCodec() codec.Codec {
	return app.appCodec
}

func (app *SimApp) ExportAppStateAndValidators(
	_ bool, _ []string, _ []string,
) (serverTypes.ExportedApp, error) {
	panic("UNIMPLEMENTED")
}

func (app *SimApp) InterfaceRegistry() codecTypes.InterfaceRegistry {
	return app.interfaceRegistry
}

func (app *SimApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *SimApp) SimulationManager() *module.SimulationManager {
	panic("UNIMPLEMENTED")
}

func (app *SimApp) TxConfig() client.TxConfig {
	return app.txConfig
}
