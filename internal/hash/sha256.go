package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type SHA256 struct {
	key [32]byte
}

func NewSHA256(keyStr string) *SHA256 {
	key := sha256.Sum256([]byte(keyStr))
	return &SHA256{key: key}
}

func (s *SHA256) Hash(data []byte) []byte {
	h := hmac.New(sha256.New, s.key[:])
	h.Write(data)
	return h.Sum(nil)
}

func (s *SHA256) StringHash(data []byte) string {
	return hex.EncodeToString(s.Hash(data))
}

func (s *SHA256) Compare(data []byte, hash []byte) bool {
	return hmac.Equal(s.Hash(data), hash)
}

func (s *SHA256) DecodeString(hashSumString string) ([]byte, error) {
	binaryHash, err := hex.DecodeString(hashSumString)
	if err != nil {

		return []byte{}, err
	}
	return binaryHash, nil
}
