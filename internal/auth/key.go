package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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
