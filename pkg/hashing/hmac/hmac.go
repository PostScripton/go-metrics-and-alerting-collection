package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/hashing"
)

type Signer struct{}

var _ hashing.Signer = &Signer{}

func NewHmacSigner() *Signer {
	return &Signer{}
}

func (s *Signer) Hash(data string, key string) []byte {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(data))
	return hash.Sum(nil)
}

func (s *Signer) HashToHex(hash []byte) string {
	return hex.EncodeToString(hash)
}

func (s *Signer) ValidHash(sign []byte, hash string) bool {
	data, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}

	return hmac.Equal(sign, data)
}
