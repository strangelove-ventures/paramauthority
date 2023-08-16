package types

import authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

const ModuleName = "authority"

var ModuleAddress = authTypes.NewModuleAddress(ModuleName)

var AuthorityKey = []byte("authority")
