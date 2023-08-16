package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	// Authority
	"github.com/noble-assets/paramauthority/x/authority/types"
	// Gov
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

var _ types.MsgServer = &Keeper{}

func (k *Keeper) Execute(goCtx context.Context, m *types.MsgExecute) (*types.MsgExecuteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	cacheCtx, writeCache := ctx.CacheContext()
	var events sdk.Events

	authority := k.GetAuthority(ctx)
	if authority != m.Authority {
		return nil, errors.ErrUnauthorized.Wrapf("expected %s, got %s", authority, m.Authority)
	}

	msgs, err := m.GetMsgs()
	if err != nil {
		return nil, err
	}

	for idx, msg := range msgs {
		handler := k.router.Handler(msg)

		res, err := handler(cacheCtx, msg)
		if err != nil {
			return nil, errors.ErrLogic.Wrapf("msg %d (%s) failed on execution: %s", idx, sdk.MsgTypeURL(msg), err)
		}

		events = append(events, res.GetEvents()...)
	}

	writeCache()
	// TODO: Events are emitted on parent ctx when written. Do we need this?
	ctx.EventManager().EmitEvents(events)

	return &types.MsgExecuteResponse{}, nil
}

func (k *Keeper) ExecuteLegacyContent(goCtx context.Context, m *types.MsgExecuteLegacyContent) (*types.MsgExecuteLegacyContentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority := types.ModuleAddress.String()
	if authority != m.Authority {
		return nil, errors.ErrUnauthorized.Wrapf("expected %s, got %s", authority, m.Authority)
	}

	content, err := v1.LegacyContentFromMessage(v1.NewMsgExecLegacyContent(m.Content, m.Authority))
	if err != nil {
		return nil, err
	}

	if !k.legacyRouter.HasRoute(content.ProposalRoute()) {
		return nil, errors.ErrLogic.Wrapf("unrecognised proposal type %s", content.ProposalType())
	}

	handler := k.legacyRouter.GetRoute(content.ProposalRoute())
	if err := handler(ctx, content); err != nil {
		return nil, errors.ErrLogic.Wrapf("proposal (%s) failed on execution: %s", content.ProposalType(), err)
	}

	return &types.MsgExecuteLegacyContentResponse{}, nil
}

func (k *Keeper) UpdateAuthority(goCtx context.Context, m *types.MsgUpdateAuthority) (*types.MsgUpdateAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority := k.GetAuthority(ctx)
	if authority != m.Authority {
		return nil, errors.ErrUnauthorized.Wrapf("expected %s, got %s", authority, m.Authority)
	}

	k.SetAuthority(ctx, m.NewAuthority)

	return &types.MsgUpdateAuthorityResponse{}, nil
}
