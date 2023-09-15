package keeper

import (
	"github.com/strangelove-ventures/paramauthority/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetParams sets the total set of upgrade parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetAuthority returns the module authority
func (k Keeper) GetAuthority(
	ctx sdk.Context,
) (authority string) {
	var res string
	k.paramSpace.Get(ctx, []byte(types.AuthorityKey), &res)
	return res
}
