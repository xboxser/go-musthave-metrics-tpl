package keypair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type PublicKey struct {
	PublicKey *x509.Certificate
}

func NewPublicKey(certPath string) (*PublicKey, error) {
	certificateBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	certificatePemBlock, _ := pem.Decode(certificateBytes)
	if certificatePemBlock == nil {
		return nil, err
	}

	certificate, err := x509.ParseCertificate(certificatePemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &PublicKey{PublicKey: certificate}, nil
}

func (p *PublicKey) Encrypt(data []byte) ([]byte, error) {
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, p.PublicKey.PublicKey.(*rsa.PublicKey), data)
	if err != nil {
		return nil, err
	}
	return encryptedMessage, nil
}
