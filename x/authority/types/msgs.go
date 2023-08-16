package types

import (
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	_, _    legacytx.LegacyMsg                 = &MsgExecute{}, &MsgUpdateAuthority{}
	_, _, _ sdk.Msg                            = &MsgExecute{}, &MsgExecuteLegacyContent{}, &MsgUpdateAuthority{}
	_, _    codecTypes.UnpackInterfacesMessage = &MsgExecute{}, &MsgExecuteLegacyContent{}
)

// MsgExecute

func (m *MsgExecute) GetMsgs() ([]sdk.Msg, error) {
	msgs, err := tx.GetMsgs(m.Messages, "")
	if err != nil {
		return nil, err
	}

	for idx, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return nil, errors.ErrUnauthorized.Wrapf("msg %d (%s) must only have one signer", idx, sdk.MsgTypeURL(msg))
		}
		if !signers[0].Equals(ModuleAddress) {
			return nil, errors.ErrUnauthorized.Wrapf("msg %d (%s) signer must be the authority address (%s)", idx, sdk.MsgTypeURL(msg), ModuleAddress)
		}

		if err := msg.ValidateBasic(); err != nil {
			return nil, errors.ErrLogic.Wrapf("msg %d (%s) failed validation: %s", idx, sdk.MsgTypeURL(msg), err)
		}
	}

	return msgs, nil
}

func (m *MsgExecute) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgExecute) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{authority}
}

func (m *MsgExecute) Route() string {
	return ModuleName
}

func (m *MsgExecute) Type() string {
	return "noble/x/authority/MsgExecute"
}

func (m *MsgExecute) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	return tx.UnpackInterfaces(unpacker, m.Messages)
}

func (m *MsgExecute) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	_, err := m.GetMsgs()
	return err
}

// MsgExecuteLegacyContent

func (m *MsgExecuteLegacyContent) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{authority}
}

func (m *MsgExecuteLegacyContent) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var content v1beta1.Content
	return unpacker.UnpackAny(m.Content, &content)
}

func (m *MsgExecuteLegacyContent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	return nil
}

// MsgUpdateAuthority

func (m *MsgUpdateAuthority) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshal(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateAuthority) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{authority}
}

func (m *MsgUpdateAuthority) Route() string {
	return ModuleName
}

func (m *MsgUpdateAuthority) Type() string {
	return "noble/x/authority/MsgUpdateAuthority"
}

func (m *MsgUpdateAuthority) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(m.NewAuthority); err != nil {
		return errors.ErrInvalidAddress.Wrapf("invalid new authority address: %s", err)
	}

	return nil
}
