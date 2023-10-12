package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
)

type ClientKeeper interface {
	ClientUpdateProposal(ctx sdk.Context, proposal *clienttypes.ClientUpdateProposal) error
	HandleUpgradeProposal(ctx sdk.Context, proposal *clienttypes.UpgradeProposal) error
}

type AuthorityKeeper interface {
	GetAuthority(ctx sdk.Context) string
}
