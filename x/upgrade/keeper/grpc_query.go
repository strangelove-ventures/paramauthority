package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	sdkupgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/strangelove-ventures/paramauthority/x/upgrade/types"
)

var _ types.QueryServer = Keeper{}

// CurrentPlan implements the Query/CurrentPlan gRPC method
func (k Keeper) CurrentPlan(c context.Context, req *sdkupgradetypes.QueryCurrentPlanRequest) (*sdkupgradetypes.QueryCurrentPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	plan, found := k.GetUpgradePlan(ctx)
	if !found {
		return &sdkupgradetypes.QueryCurrentPlanResponse{}, nil
	}

	return &sdkupgradetypes.QueryCurrentPlanResponse{Plan: &plan}, nil
}

// AppliedPlan implements the Query/AppliedPlan gRPC method
func (k Keeper) AppliedPlan(c context.Context, req *sdkupgradetypes.QueryAppliedPlanRequest) (*sdkupgradetypes.QueryAppliedPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	applied := k.GetDoneHeight(ctx, req.Name)
	if applied == 0 {
		return &sdkupgradetypes.QueryAppliedPlanResponse{}, nil
	}

	return &sdkupgradetypes.QueryAppliedPlanResponse{Height: applied}, nil
}

// UpgradedConsensusState implements the Query/UpgradedConsensusState gRPC method
// nolint: staticcheck
func (k Keeper) UpgradedConsensusState(c context.Context, req *sdkupgradetypes.QueryUpgradedConsensusStateRequest) (*sdkupgradetypes.QueryUpgradedConsensusStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	consState, found := k.GetUpgradedConsensusState(ctx, req.LastHeight)
	if !found {
		return &sdkupgradetypes.QueryUpgradedConsensusStateResponse{}, nil
	}

	return &sdkupgradetypes.QueryUpgradedConsensusStateResponse{
		UpgradedConsensusState: consState,
	}, nil
}

// ModuleVersions implements the Query/QueryModuleVersions gRPC method
func (k Keeper) ModuleVersions(c context.Context, req *sdkupgradetypes.QueryModuleVersionsRequest) (*sdkupgradetypes.QueryModuleVersionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// check if a specific module was requested
	if len(req.ModuleName) > 0 {
		if version, ok := k.getModuleVersion(ctx, req.ModuleName); ok {
			// return the requested module
			res := []*sdkupgradetypes.ModuleVersion{{Name: req.ModuleName, Version: version}}
			return &sdkupgradetypes.QueryModuleVersionsResponse{ModuleVersions: res}, nil
		}
		// module requested, but not found
		return nil, errors.Wrapf(errors.ErrNotFound, "x/upgrade: QueryModuleVersions module %s not found", req.ModuleName)
	}

	// if no module requested return all module versions from state
	mv := k.GetModuleVersions(ctx)
	return &sdkupgradetypes.QueryModuleVersionsResponse{
		ModuleVersions: mv,
	}, nil
}

// Authority implements the Query/Authority gRPC method, returning the account capable of performing upgrades
func (k Keeper) Authority(c context.Context, req *types.QueryAuthorityRequest) (*types.QueryAuthorityResponse, error) {
	return &types.QueryAuthorityResponse{Address: k.authority}, nil
}
