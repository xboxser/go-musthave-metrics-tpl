package key_pair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type PrivateKey struct {
	PrivateKey *rsa.PrivateKey
}

func NewPrivateKey(certPath string) (*PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	privateKeyPemBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyPemBlock == nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyPemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{PrivateKey: privateKey}, nil
}

func (p *PrivateKey) Decrypt(data []byte) ([]byte, error) {
	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, p.PrivateKey, data)
	if err != nil {
		return nil, err
	}
	return decryptedMessage, nil
}
