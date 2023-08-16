package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(authority string) *GenesisState {
	return &GenesisState{
		Authority: authority,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState("")
}

// Validate validates the provided genesis state.
func (gs *GenesisState) Validate() error {
	if gs.Authority != "" {
		if _, err := sdk.AccAddressFromBech32(gs.Authority); err != nil {
			return err
		}
	}

	return nil
}
