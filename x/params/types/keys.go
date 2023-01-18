package types

const (
	// ModuleName defines the module name
	ModuleName = "params"

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	AuthorityKey = "authority"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
