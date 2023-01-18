package proposal

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the params module
func NewParams(authority string) Params {
	return Params{
		Authority: authority,
	}
}

// DefaultParams is the default parameter configuration for the params module
func DefaultParams() Params {
	return NewParams("")
}

// Validate all params module parameters
func (p Params) Validate() error {
	return validateAuthority(p.Authority)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte(AuthorityKey), p.Authority, validateAuthority),
	}
}

func validateAuthority(i interface{}) error {
	a, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if _, err := sdk.AccAddressFromBech32(a); err != nil {
		return sdkerrors.Wrap(err, "authority")
	}

	return nil
}
