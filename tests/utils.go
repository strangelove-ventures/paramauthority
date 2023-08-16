package tests

import (
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

func RandomAccount(cfg ibc.ChainConfig) ibc.Wallet {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	path := hd.CreateHDPath(118, 0, 0)
	rawPk, _ := hd.Secp256k1.Derive()(mnemonic, "", path.String())

	pk := secp256k1.PrivKey(rawPk)
	return cosmos.NewWallet("authority", pk.PubKey().Address(), mnemonic, cfg)
}
