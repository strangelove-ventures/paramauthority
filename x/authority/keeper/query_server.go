package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/paramauthority/x/authority/types"
)

var _ types.QueryServer = &Keeper{}

func (k *Keeper) Authority(goCtx context.Context, _ *types.QueryAuthority) (*types.QueryAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	authority := k.GetAuthority(ctx)

	return &types.QueryAuthorityResponse{Authority: authority}, nil
}
