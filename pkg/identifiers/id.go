package identifiers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func generate(length int) (string, error) {
	if length%2 != 0 {
		return "", errors.New("length must be even")
	}
	buf := make([]byte, length/2)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

const (
	IDLength      = 16
	ShortIDLength = 8
)

func GenerateID() (string, error) {
	return generate(IDLength)
}

func GenerateShortID() (string, error) {
	return generate(ShortIDLength)
}
