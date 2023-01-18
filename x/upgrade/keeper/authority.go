package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/paramauthority/x/params/types"
)

// SetAuthority sets the module authority
func (k Keeper) SetAuthority(ctx sdk.Context, authority string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.AuthorityKey), []byte(authority))
}

// GetAuthority returns the module authority
func (k Keeper) GetAuthority(
	ctx sdk.Context,
) (authority string, found bool) {
	store := ctx.KVStore(k.storeKey)
	a := store.Get(types.KeyPrefix(types.AuthorityKey))
	if a == nil {
		return "", false
	}
	return string(a), true
}
