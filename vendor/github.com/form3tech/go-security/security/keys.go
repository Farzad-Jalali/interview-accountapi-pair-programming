package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

type TestKeyPair struct {
	PublicKeyPem  string
	PublicKeyDer  string
	PrivateKeyPem string
	RsaPrivateKey *rsa.PrivateKey
}

func GenerateTestKeyPair() (*TestKeyPair, error) {

	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	publicKey := &key.PublicKey

	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)

	if err != nil {
		return nil, fmt.Errorf("could not marshall public key, %v", err)

	}

	pemPublicKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}))

	var privateKeyPem = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return &TestKeyPair{
		PublicKeyPem:  pemPublicKey,
		PublicKeyDer:  base64.StdEncoding.EncodeToString(publicKeyBytes),
		PrivateKeyPem: string(pem.EncodeToMemory(privateKeyPem)),
		RsaPrivateKey: key,
	}, nil

}
