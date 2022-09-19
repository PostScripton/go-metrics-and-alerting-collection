package hashing

type Signer interface {
	Hash(data string, key string) []byte
	HashToHex(hash []byte) string
	ValidHash(sign []byte, hash string) bool
}
