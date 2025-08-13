package hash

type Hasher interface {
	Hash(data []byte) []byte
	StringHash(data []byte) string
	Compare(data []byte, hash []byte) bool
	DecodeString(hashSumString string) ([]byte, error)
}
