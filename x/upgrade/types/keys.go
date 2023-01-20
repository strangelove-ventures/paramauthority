package types

const (
	AuthorityKey = "authority"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
