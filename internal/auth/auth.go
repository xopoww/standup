package auth

import (
	"crypto/ecdsa"
	"errors"
)

type Authenticator interface {
	Enabled() bool
	Authenticate(raw string) (subjectId string, err error)
}

type staticAuthenticator struct {
	enabled bool
	key     *ecdsa.PublicKey
}

func (a staticAuthenticator) Enabled() bool {
	return a.enabled
}

func (a staticAuthenticator) Authenticate(raw string) (subjectId string, err error) {
	if !a.enabled {
		return "", errors.New("authentication is disabled")
	}
	return ValidateToken(raw, a.key)
}

func NewStaticAuthenticator(key *ecdsa.PublicKey) Authenticator {
	return &staticAuthenticator{
		enabled: true,
		key:     key,
	}
}

func NewDisabledAuthenticator() Authenticator {
	return &staticAuthenticator{enabled: false}
}
