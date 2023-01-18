package proposal

const (
	// RouterKey defines the routing key for a ParameterChangeProposal
	RouterKey = "params"

	AuthorityKey = "authority"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
