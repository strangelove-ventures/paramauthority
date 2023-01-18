package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the global paramstore
type Keeper struct {
	keeper.Keeper
	storeKey sdk.StoreKey
}

// NewKeeper constructs a params keeper
func NewKeeper(cdc codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) Keeper {
	return Keeper{
		Keeper: keeper.NewKeeper(
			cdc,
			legacyAmino,
			key,
			tkey,
		),
		storeKey: key,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return k.Keeper.Logger(ctx)
}

// Allocate subspace used for keepers
func (k Keeper) Subspace(s string) types.Subspace {
	return k.Keeper.Subspace(s)
}

// Get existing substore from keeper
func (k Keeper) GetSubspace(s string) (types.Subspace, bool) {
	return k.Keeper.GetSubspace(s)
}
