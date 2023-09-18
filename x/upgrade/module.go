package upgrade

import (
	"context"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdkupgrademodule "github.com/cosmos/cosmos-sdk/x/upgrade"
	sdkupgradecli "github.com/cosmos/cosmos-sdk/x/upgrade/client/cli"
	sdkupgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	sdkupgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/strangelove-ventures/paramauthority/x/upgrade/client/cli"
	"github.com/strangelove-ventures/paramauthority/x/upgrade/keeper"
	"github.com/strangelove-ventures/paramauthority/x/upgrade/types"
)

func init() {
	sdkupgradetypes.RegisterLegacyAminoCodec(codec.NewLegacyAmino())
}

const (
	consensusVersion uint64 = 1
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic implements the sdk.AppModuleBasic interface
type AppModuleBasic struct{}

// Name returns the ModuleName
func (AppModuleBasic) Name() string {
	return sdkupgradetypes.ModuleName
}

// RegisterLegacyAminoCodec registers the upgrade types on the LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
	sdkupgradetypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterRESTRoutes registers all REST query handlers
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, r *mux.Router) {

}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the upgrade module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := sdkupgradetypes.RegisterQueryHandlerClient(context.Background(), mux, sdkupgradetypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetQueryCmd returns the cli query commands for this module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return sdkupgradecli.GetQueryCmd()
}

// GetTxCmd returns the transaction commands for this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
	sdkupgradetypes.RegisterInterfaces(registry)
}

// AppModule implements the sdk.AppModule interface
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// RegisterInvariants does nothing, there are no invariants to enforce
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Deprecated: Route returns the message routing key for the upgrade module.
func (AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute returns the route we respond to for abci queries
func (AppModule) QuerierRoute() string { return sdkupgradetypes.QuerierKey }

// LegacyQuerierHandler registers a query handler to respond to the module-specific queries
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return sdkupgradekeeper.NewQuerier(am.keeper.Keeper, legacyQuerierCdc)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	sdkupgradetypes.RegisterQueryServer(cfg.QueryServer(), am.keeper.Keeper)
}

// InitGenesis is ignored, no sense in serializing future upgrades
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	am.keeper.SetParams(ctx, genState.Params)

	return []abci.ValidatorUpdate{}
}

// DefaultGenesis is an empty object
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis is always successful, as we ignore the value
func (AppModuleBasic) ValidateGenesis(_ codec.JSONCodec, config client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}

// ExportGenesis is always empty, as InitGenesis does nothing either
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	authority := am.keeper.GetAuthority(ctx)

	genState := &types.GenesisState{
		Params: types.NewParams(authority),
	}

	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return consensusVersion }

// BeginBlock calls the upgrade module hooks
//
// CONTRACT: this is registered in BeginBlocker *before* all other modules' BeginBlock functions
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	sdkupgrademodule.BeginBlocker(am.keeper.Keeper, ctx, req)
}
