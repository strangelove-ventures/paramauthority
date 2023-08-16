package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/paramauthority/x/authority/types"
)

// GetAuthority returns the authority from state.
func (k *Keeper) GetAuthority(ctx sdk.Context) (authority string) {
	bz := ctx.KVStore(k.storeKey).Get(types.AuthorityKey)
	if bz == nil {
		panic("authority not found in state")
	}

	return string(bz)
}

// SetAuthority stores the authority in state.
func (k *Keeper) SetAuthority(ctx sdk.Context, authority string) {
	bz := []byte(authority)
	ctx.KVStore(k.storeKey).Set(types.AuthorityKey, bz)
}
