package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/paramauthority/x/upgrade/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the upgrade MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

var _ types.MsgServer = msgServer{}

// SoftwareUpgrade implements the Msg/SoftwareUpgrade Msg service.
func (k msgServer) SoftwareUpgrade(goCtx context.Context, req *types.MsgSoftwareUpgrade) (*types.MsgSoftwareUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority := k.Keeper.GetAuthority(ctx)
	if err := types.NewParams(authority).Validate(); err != nil {
		return nil, err
	}

	if authority != req.Authority {
		return nil, fmt.Errorf("invalid signer: expected %s got %s", authority, req.Authority)
	}

	err := k.ScheduleUpgrade(ctx, req.Plan)
	if err != nil {
		return nil, err
	}

	return &types.MsgSoftwareUpgradeResponse{}, nil
}

// CancelUpgrade implements the Msg/CancelUpgrade Msg service.
func (k msgServer) CancelUpgrade(goCtx context.Context, req *types.MsgCancelUpgrade) (*types.MsgCancelUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority := k.Keeper.GetAuthority(ctx)
	if err := types.NewParams(authority).Validate(); err != nil {
		return nil, err
	}

	if authority != req.Authority {
		return nil, fmt.Errorf("invalid signer: expected %s got %s", authority, req.Authority)
	}

	k.ClearUpgradePlan(ctx)

	return &types.MsgCancelUpgradeResponse{}, nil
}
