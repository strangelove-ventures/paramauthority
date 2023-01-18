package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	xp "github.com/cosmos/cosmos-sdk/x/upgrade/exported"
	sdkupgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	sdkupgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Deprecated: UpgradeInfoFileName file to store upgrade information
// use x/upgrade/types.UpgradeInfoFilename instead.
const UpgradeInfoFileName string = "upgrade-info.json"

type Keeper struct {
	sdkupgradekeeper.Keeper
	storeKey storetypes.StoreKey // key to access x/upgrade store
}

// NewKeeper constructs an upgrade Keeper which requires the following arguments:
// skipUpgradeHeights - map of heights to skip an upgrade
// storeKey - a store key with which to access upgrade's store
// cdc - the app-wide binary codec
// homePath - root directory of the application's config
// vs - the interface implemented by baseapp which allows setting baseapp's protocol version field
func NewKeeper(skipUpgradeHeights map[int64]bool, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, homePath string, vs xp.ProtocolVersionSetter) Keeper {
	return Keeper{
		Keeper:   sdkupgradekeeper.NewKeeper(skipUpgradeHeights, storeKey, cdc, homePath, vs),
		storeKey: storeKey,
	}
}

// SetUpgradeHandler sets an UpgradeHandler for the upgrade specified by name. This handler will be called when the upgrade
// with this name is applied. In order for an upgrade with the given name to proceed, a handler for this upgrade
// must be set even if it is a no-op function.
func (k Keeper) SetUpgradeHandler(name string, upgradeHandler sdkupgradetypes.UpgradeHandler) {
	k.Keeper.SetUpgradeHandler(name, upgradeHandler)
}

// SetModuleVersionMap saves a given version map to state
func (k Keeper) SetModuleVersionMap(ctx sdk.Context, vm module.VersionMap) {
	k.Keeper.SetModuleVersionMap(ctx, vm)
}

// GetModuleVersionMap returns a map of key module name and value module consensus version
// as defined in ADR-041.
func (k Keeper) GetModuleVersionMap(ctx sdk.Context) module.VersionMap {
	return k.Keeper.GetModuleVersionMap(ctx)
}

// GetModuleVersions gets a slice of module consensus versions
func (k Keeper) GetModuleVersions(ctx sdk.Context) []*sdkupgradetypes.ModuleVersion {
	return k.Keeper.GetModuleVersions(ctx)
}

// ScheduleUpgrade schedules an upgrade based on the specified plan.
// If there is another Plan already scheduled, it will cancel and overwrite it.
// ScheduleUpgrade will also write the upgraded IBC ClientState to the upgraded client
// path if it is specified in the plan.
func (k Keeper) ScheduleUpgrade(ctx sdk.Context, plan sdkupgradetypes.Plan) error {
	return k.Keeper.ScheduleUpgrade(ctx, plan)
}

// SetUpgradedClient sets the expected upgraded client for the next version of this chain at the last height the current chain will commit.
func (k Keeper) SetUpgradedClient(ctx sdk.Context, planHeight int64, bz []byte) error {
	return k.Keeper.SetUpgradedClient(ctx, planHeight, bz)
}

// GetUpgradedClient gets the expected upgraded client for the next version of this chain
func (k Keeper) GetUpgradedClient(ctx sdk.Context, height int64) ([]byte, bool) {
	return k.Keeper.GetUpgradedClient(ctx, height)
}

// SetUpgradedConsensusState set the expected upgraded consensus state for the next version of this chain
// using the last height committed on this chain.
func (k Keeper) SetUpgradedConsensusState(ctx sdk.Context, planHeight int64, bz []byte) error {
	return k.Keeper.SetUpgradedConsensusState(ctx, planHeight, bz)
}

// GetUpgradedConsensusState set the expected upgraded consensus state for the next version of this chain
func (k Keeper) GetUpgradedConsensusState(ctx sdk.Context, lastHeight int64) ([]byte, bool) {
	return k.Keeper.GetUpgradedConsensusState(ctx, lastHeight)
}

// GetLastCompletedUpgrade returns the last applied upgrade name and height.
func (k Keeper) GetLastCompletedUpgrade(ctx sdk.Context) (string, int64) {
	return k.Keeper.GetLastCompletedUpgrade(ctx)
}

// GetDoneHeight returns the height at which the given upgrade was executed
func (k Keeper) GetDoneHeight(ctx sdk.Context, name string) int64 {
	return k.Keeper.GetDoneHeight(ctx, name)
}

// ClearIBCState clears any planned IBC state
func (k Keeper) ClearIBCState(ctx sdk.Context, lastHeight int64) {
	k.Keeper.ClearIBCState(ctx, lastHeight)
}

// ClearUpgradePlan clears any schedule upgrade and associated IBC states.
func (k Keeper) ClearUpgradePlan(ctx sdk.Context) {
	k.Keeper.ClearUpgradePlan(ctx)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return k.Keeper.Logger(ctx)
}

// GetUpgradePlan returns the currently scheduled Plan if any, setting havePlan to true if there is a scheduled
// upgrade or false if there is none
func (k Keeper) GetUpgradePlan(ctx sdk.Context) (plan sdkupgradetypes.Plan, havePlan bool) {
	return k.Keeper.GetUpgradePlan(ctx)
}

// HasHandler returns true iff there is a handler registered for this name
func (k Keeper) HasHandler(name string) bool {
	return k.Keeper.HasHandler(name)
}

// ApplyUpgrade will execute the handler associated with the Plan and mark the plan as done.
func (k Keeper) ApplyUpgrade(ctx sdk.Context, plan sdkupgradetypes.Plan) {
	k.Keeper.ApplyUpgrade(ctx, plan)
}

// IsSkipHeight checks if the given height is part of skipUpgradeHeights
func (k Keeper) IsSkipHeight(height int64) bool {
	return k.Keeper.IsSkipHeight(height)
}

// DumpUpgradeInfoToDisk writes upgrade information to UpgradeInfoFileName.
func (k Keeper) DumpUpgradeInfoToDisk(height int64, name string) error {
	return k.Keeper.DumpUpgradeInfoToDisk(height, name)
}

// GetUpgradeInfoPath returns the upgrade info file path
func (k Keeper) GetUpgradeInfoPath() (string, error) {
	return k.Keeper.GetUpgradeInfoPath()
}

// ReadUpgradeInfoFromDisk returns the name and height of the upgrade which is
// written to disk by the old binary when panicking. An error is returned if
// the upgrade path directory cannot be created or if the file exists and
// cannot be read or if the upgrade info fails to unmarshal.
func (k Keeper) ReadUpgradeInfoFromDisk() (storetypes.UpgradeInfo, error) {
	return k.Keeper.ReadUpgradeInfoFromDisk()
}

// SetDowngradeVerified updates downgradeVerified.
func (k *Keeper) SetDowngradeVerified(v bool) {
	k.Keeper.SetDowngradeVerified(v)
}

// DowngradeVerified returns downgradeVerified.
func (k Keeper) DowngradeVerified() bool {
	return k.Keeper.DowngradeVerified()
}
