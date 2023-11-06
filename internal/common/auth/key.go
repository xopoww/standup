package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

const (
	privateKey = "EC PRIVATE KEY"
	publicKey  = "PUBLIC KEY"
)

func loadKey(path, blockType string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid pem")
	}
	if block.Type != blockType {
		return nil, fmt.Errorf("wrong block type: want %q, got %q", blockType, block.Type)
	}
	return block.Bytes, nil
}

func writeKey(data []byte, blockType string, w io.Writer) error {
	err := pem.Encode(w, &pem.Block{
		Type:  blockType,
		Bytes: data,
	})
	if err != nil {
		return fmt.Errorf("pem encode: %w", err)
	}
	return nil
}

func LoadPublicKey(path string) (*ecdsa.PublicKey, error) {
	encoded, err := loadKey(path, publicKey)
	if err != nil {
		return nil, fmt.Errorf("load pem: %w", err)
	}
	genericKey, err := x509.ParsePKIXPublicKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("x509 parse: %w", err)
	}
	key, ok := genericKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected key algo")
	}
	return key, nil
}

func WritePublicKey(key *ecdsa.PublicKey, w io.Writer) error {
	encoded, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("x509 marshal: %w", err)
	}
	err = writeKey(encoded, publicKey, w)
	if err != nil {
		return fmt.Errorf("write key: %w", err)
	}
	return nil
}

func LoadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	encoded, err := loadKey(path, privateKey)
	if err != nil {
		return nil, fmt.Errorf("load pem: %w", err)
	}
	key, err := x509.ParseECPrivateKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("x509 parse: %w", err)
	}
	return key, nil
}

func WritePrivateKey(key *ecdsa.PrivateKey, w io.Writer) error {
	encoded, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return fmt.Errorf("x509 marshal: %w", err)
	}
	err = writeKey(encoded, privateKey, w)
	if err != nil {
		return fmt.Errorf("write key: %w", err)
	}
	return nil
}
