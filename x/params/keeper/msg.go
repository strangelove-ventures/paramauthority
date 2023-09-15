package keeper

import (
	"context"
	"fmt"

	"github.com/strangelove-ventures/paramauthority/x/params/types/proposal"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the upgrade MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) proposal.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

var _ proposal.MsgServer = msgServer{}

// SoftwareUpgrade implements the Msg/SoftwareUpgrade Msg service.
func (k msgServer) UpdateParams(goCtx context.Context, req *proposal.MsgUpdateParams) (*proposal.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority := k.Keeper.GetAuthority(ctx)
	if err := proposal.NewParams(authority).Validate(); err != nil {
		return nil, err
	}

	if authority != req.Authority {
		return nil, fmt.Errorf("expected %s got %s", authority, req.Authority)
	}

	for _, c := range req.ChangeProposal.Changes {
		ss, ok := k.Keeper.Keeper.GetSubspace(c.Subspace)
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

	return &proposal.MsgUpdateParamsResponse{}, nil
}
