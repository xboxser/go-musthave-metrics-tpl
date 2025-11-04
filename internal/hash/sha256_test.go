package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	type args struct {
		contentType string
	}

	tests := []struct {
		name string
		key  string
		data []byte
	}{
		{name: "simple test case", key: "test_key", data: []byte("test_data")},
		{name: "empty data", key: "test_key", data: []byte("")},
		{name: "empty key", key: "", data: []byte("tester")},
		{name: "empty all", key: "", data: []byte("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hasher := NewSHA256(tt.key)
			// Вычисляем хэш
			actualHash := hasher.Hash(tt.data)

			// Вычисляем ожидаемый хэш вручную для проверки
			expectedKey := sha256.Sum256([]byte(tt.key))
			expectedHasher := hmac.New(sha256.New, expectedKey[:])
			expectedHasher.Write(tt.data)
			expectedHash := expectedHasher.Sum(nil)

			if len(actualHash) != len(expectedHash) {
				t.Errorf("hash lengths don't match: actual %d, expected %d", len(actualHash), len(expectedHash))
			}
			for i := range actualHash {
				if actualHash[i] != expectedHash[i] {
					t.Errorf("hash byte at index %d doesn't match: %d, expected %d", i, actualHash[i], expectedHash[i])
				}
			}
		})
	}
}

func TestStringHash(t *testing.T) {
	type args struct {
		contentType string
	}

	tests := []struct {
		name string
		key  string
		data []byte
	}{
		{name: "simple test case", key: "test_key", data: []byte("test_data")},
		{name: "empty data", key: "test_key", data: []byte("")},
		{name: "empty key", key: "", data: []byte("tester")},
		{name: "empty all", key: "", data: []byte("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewSHA256(tt.key)
			// Вычисляем хэш
			actualHash := hasher.StringHash(tt.data)

			// Вычисляем ожидаемый хэш вручную для проверки
			expectedKey := sha256.Sum256([]byte(tt.key))
			expectedHasher := hmac.New(sha256.New, expectedKey[:])
			expectedHasher.Write(tt.data)
			expectedHash := hex.EncodeToString(expectedHasher.Sum(nil))

			if len(actualHash) != len(expectedHash) {
				t.Errorf("hash lengths don't match: actual %d, expected %d", len(actualHash), len(expectedHash))
			}
			if actualHash != expectedHash {
				t.Errorf("hashes don't match: %s, expected %s", actualHash, expectedHash)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	type args struct {
		contentType string
	}

	tests := []struct {
		name string
		key  string
		data []byte
	}{
		{name: "simple test case", key: "test_key", data: []byte("test_data")},
		{name: "empty data", key: "test_key", data: []byte("")},
		{name: "empty key", key: "", data: []byte("tester")},
		{name: "empty all", key: "", data: []byte("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewSHA256(tt.key)

			// Вычисляем ожидаемый хэш вручную для проверки
			expectedKey := sha256.Sum256([]byte(tt.key))
			expectedHasher := hmac.New(sha256.New, expectedKey[:])
			expectedHasher.Write(tt.data)
			expectedHash := expectedHasher.Sum(nil)

			if !hasher.Compare(tt.data, expectedHash) {
				t.Errorf("incorrect comparison name: %s", tt.name)
			}
		})
	}
}
