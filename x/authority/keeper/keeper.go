package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storeTypes.StoreKey

		router       *baseapp.MsgServiceRouter
		legacyRouter v1beta1.Router
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storeTypes.StoreKey,
	router *baseapp.MsgServiceRouter,
	legacyRouter v1beta1.Router,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		router:       router,
		legacyRouter: legacyRouter,
	}
}

func (k *Keeper) SetLegacyRouter(legacyRouter v1beta1.Router) {
	k.legacyRouter = legacyRouter
}
