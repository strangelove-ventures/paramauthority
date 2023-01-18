package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/strangelove-ventures/paramauthority/x/params/types/proposal"
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
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, found := k.Keeper.GetAuthority(ctx)
	if !found {
		return nil, fmt.Errorf("authority does not exist")
	}

	if authority != req.Authority {
		return nil, fmt.Errorf("expected %s got %s", authority, req.Authority)
	}

	for _, c := range req.ChangeProposal.Changes {
		ss, ok := k.GetSubspace(c.Subspace)
		if !ok {
			return nil, fmt.Errorf("unknown subspace, %s", c.Subspace)
		}

		k.Logger(ctx).Info(
			fmt.Sprintf("attempt to set new parameter value; key: %s, value: %s", c.Key, c.Value),
		)

		if err := ss.Update(ctx, []byte(c.Key), []byte(c.Value)); err != nil {
			return nil, fmt.Errorf("key: %s, value: %s, err: %w", c.Key, c.Value, err)
		}
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
