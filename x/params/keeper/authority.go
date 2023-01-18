package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/paramauthority/x/params/types"
	"github.com/strangelove-ventures/paramauthority/x/params/types/proposal"
)

// SetParams sets the total set of params parameters.
func (k Keeper) SetParams(ctx sdk.Context, params proposal.Params) {
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
