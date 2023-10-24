package ibc

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	"github.com/strangelove-ventures/paramauthority/x/ibc/types"
	upgradetypes "github.com/strangelove-ventures/paramauthority/x/upgrade/types"
)

type msgServer struct {
	authorityKeeper types.AuthorityKeeper
	clientKeeper    types.ClientKeeper
}

func NewMsgServer(
	authorityKeeper types.AuthorityKeeper,
	clientKeeper types.ClientKeeper,
) types.MsgServer {
	return msgServer{
		authorityKeeper: authorityKeeper,
		clientKeeper:    clientKeeper,
	}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) ClientUpdate(goCtx context.Context, req *types.MsgClientUpdate) (*types.MsgClientUpdateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority := m.authorityKeeper.GetAuthority(ctx)
	if err := upgradetypes.NewParams(authority).Validate(); err != nil {
		return nil, err
	}

	if authority != req.Authority {
		return nil, fmt.Errorf("invalid signer: expected %s got %s", authority, req.Authority)
	}

	if err := m.clientKeeper.ClientUpdateProposal(ctx, &clienttypes.ClientUpdateProposal{
		SubjectClientId:    req.SubjectClientId,
		SubstituteClientId: req.SubstituteClientId,
	}); err != nil {
		return nil, err
	}

	return &types.MsgClientUpdateResponse{}, nil
}

func (m msgServer) Upgrade(goCtx context.Context, req *types.MsgUpgrade) (*types.MsgUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority := m.authorityKeeper.GetAuthority(ctx)
	if err := upgradetypes.NewParams(authority).Validate(); err != nil {
		return nil, err
	}

	if authority != req.Authority {
		return nil, fmt.Errorf("invalid signer: expected %s got %s", authority, req.Authority)
	}

	if err := m.clientKeeper.HandleUpgradeProposal(ctx, &clienttypes.UpgradeProposal{
		Plan:                req.Plan,
		UpgradedClientState: req.UpgradedClientState,
	}); err != nil {
		return nil, err
	}

	return &types.MsgUpgradeResponse{}, nil
}
