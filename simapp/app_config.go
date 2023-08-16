package simapp

import (
	runtime "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	app "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/appconfig"

	// Auth
	auth "cosmossdk.io/api/cosmos/auth/module/v1"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Authority
	authority "github.com/noble-assets/paramauthority/pulsar/noble/authority/module/v1"
	authorityTypes "github.com/noble-assets/paramauthority/x/authority/types"
	// Bank
	bank "cosmossdk.io/api/cosmos/bank/module/v1"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// Capability
	capability "cosmossdk.io/api/cosmos/capability/module/v1"
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	// Consensus
	consensus "cosmossdk.io/api/cosmos/consensus/module/v1"
	consensusTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	// GenUtil
	genUtil "cosmossdk.io/api/cosmos/genutil/module/v1"
	genUtilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	// IBC Core
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/core/exported"
	// IBC Fee
	ibcFeeTypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	// IBC Transfer
	ibcTransferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	// ICA
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	// Params
	params "cosmossdk.io/api/cosmos/params/module/v1"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	// Staking
	staking "cosmossdk.io/api/cosmos/staking/module/v1"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// Tx
	tx "cosmossdk.io/api/cosmos/tx/config/v1"
	// Upgrade
	upgrade "cosmossdk.io/api/cosmos/upgrade/module/v1"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

var AppConfig = appconfig.Compose(&app.Config{
	Modules: []*app.ModuleConfig{
		{
			Name: "runtime",
			Config: appconfig.WrapAny(&runtime.Module{
				AppName: "SimApp",
				BeginBlockers: []string{
					upgradeTypes.ModuleName,
					capabilityTypes.ModuleName,
					stakingTypes.ModuleName,
					authTypes.ModuleName,
					bankTypes.ModuleName,
					genUtilTypes.ModuleName,
					paramsTypes.ModuleName,
					consensusTypes.ModuleName,

					ibcTransferTypes.ModuleName,
					ibcTypes.ModuleName,
					icaTypes.ModuleName,
					ibcFeeTypes.ModuleName,

					authorityTypes.ModuleName,
				},
				EndBlockers: []string{
					stakingTypes.ModuleName,
					capabilityTypes.ModuleName,
					authTypes.ModuleName,
					bankTypes.ModuleName,
					genUtilTypes.ModuleName,
					paramsTypes.ModuleName,
					consensusTypes.ModuleName,
					upgradeTypes.ModuleName,

					ibcTransferTypes.ModuleName,
					ibcTypes.ModuleName,
					icaTypes.ModuleName,
					ibcFeeTypes.ModuleName,

					authorityTypes.ModuleName,
				},
				OverrideStoreKeys: []*runtime.StoreKeyConfig{
					{
						ModuleName: authTypes.ModuleName,
						KvStoreKey: "acc",
					},
				},
				InitGenesis: []string{
					capabilityTypes.ModuleName,
					authTypes.ModuleName,
					bankTypes.ModuleName,
					stakingTypes.ModuleName,
					genUtilTypes.ModuleName,
					paramsTypes.ModuleName,
					upgradeTypes.ModuleName,
					consensusTypes.ModuleName,

					ibcTransferTypes.ModuleName,
					ibcTypes.ModuleName,
					icaTypes.ModuleName,
					ibcFeeTypes.ModuleName,

					authorityTypes.ModuleName,
				},
			}),
		},

		// Cosmos SDK Modules

		{
			Name: authTypes.ModuleName,
			Config: appconfig.WrapAny(&auth.Module{
				Authority:    "authority",
				Bech32Prefix: "cosmos",
				ModuleAccountPermissions: []*auth.ModuleAccountPermission{
					{Account: authTypes.FeeCollectorName},
					{Account: stakingTypes.BondedPoolName, Permissions: []string{authTypes.Burner, stakingTypes.ModuleName}},
					{Account: stakingTypes.NotBondedPoolName, Permissions: []string{authTypes.Burner, stakingTypes.ModuleName}},

					{Account: ibcFeeTypes.ModuleName},
					{Account: ibcTransferTypes.ModuleName, Permissions: []string{authTypes.Burner, authTypes.Minter}},
					{Account: icaTypes.ModuleName},
				},
			}),
		},
		{
			Name: bankTypes.ModuleName,
			Config: appconfig.WrapAny(&bank.Module{
				Authority: "authority",
				BlockedModuleAccountsOverride: []string{
					authTypes.FeeCollectorName,
					stakingTypes.BondedPoolName,
					stakingTypes.NotBondedPoolName,

					authorityTypes.ModuleName,
				},
			}),
		},
		{
			Name: capabilityTypes.ModuleName,
			Config: appconfig.WrapAny(&capability.Module{
				SealKeeper: true,
			}),
		},
		{
			Name: consensusTypes.ModuleName,
			Config: appconfig.WrapAny(&consensus.Module{
				Authority: "authority",
			}),
		},
		{
			Name:   genUtilTypes.ModuleName,
			Config: appconfig.WrapAny(&genUtil.Module{}),
		},
		{
			Name:   paramsTypes.ModuleName,
			Config: appconfig.WrapAny(&params.Module{}),
		},
		{
			Name: stakingTypes.ModuleName,
			Config: appconfig.WrapAny(&staking.Module{
				Authority: "authority",
			}),
		},
		{
			Name:   "tx",
			Config: appconfig.WrapAny(&tx.Config{}),
		},
		{
			Name: upgradeTypes.ModuleName,
			Config: appconfig.WrapAny(&upgrade.Module{
				Authority: "authority",
			}),
		},

		// Custom Modules

		{
			Name:   authorityTypes.ModuleName,
			Config: appconfig.WrapAny(&authority.Module{}),
		},
	},
})
