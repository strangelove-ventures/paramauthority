package authority

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/paramauthority/x/authority/keeper"
	"github.com/noble-assets/paramauthority/x/authority/types"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, gs types.GenesisState) {
	k.SetAuthority(ctx, gs.Authority)
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	authority := k.GetAuthority(ctx)
	return types.NewGenesisState(authority)
}
