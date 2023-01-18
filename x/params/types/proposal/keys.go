package proposal

const (
	// ModuleName defines the name of the module
	ModuleName = "params"

	// RouterKey defines the routing key for a ParameterChangeProposal
	RouterKey = "params"

	AuthorityKey = "authority"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
