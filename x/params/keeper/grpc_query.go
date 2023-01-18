package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

var _ proposal.QueryServer = Keeper{}

// Params returns subspace params
func (k Keeper) Params(c context.Context, req *proposal.QueryParamsRequest) (*proposal.QueryParamsResponse, error) {
	return k.Keeper.Params(c, req)
}
