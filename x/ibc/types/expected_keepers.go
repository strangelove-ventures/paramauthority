package types

import (
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ClientKeeper interface {
	ClientUpdateProposal(ctx sdk.Context, proposal *clienttypes.ClientUpdateProposal) error
	HandleUpgradeProposal(ctx sdk.Context, proposal *clienttypes.UpgradeProposal) error
}

type AuthorityKeeper interface {
	GetAuthority(ctx sdk.Context) string
}
