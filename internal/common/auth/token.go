package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var signingMethod = jwt.SigningMethodES256

func IssueToken(subjectID int64, notBefore, expiresAt time.Time, privateKey any) (string, error) {
	now := time.Now()
	c := &jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   idToString(subjectID),
		Audience:  []string{},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(notBefore),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        "",
	}
	token := jwt.NewWithClaims(signingMethod, c)
	return token.SignedString(privateKey)
}

func ValidateToken(raw string, publicKey any) (subjectID int64, err error) {
	token, err := jwt.ParseWithClaims(raw, &jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) { return publicKey, nil },
		jwt.WithValidMethods([]string{signingMethod.Name}),
		jwt.WithoutClaimsValidation(),
	)
	if err != nil {
		return 0, fmt.Errorf("parse jwt: %w", err)
	}

	c, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	now := time.Now()
	if c.NotBefore == nil || now.Before(c.NotBefore.Time) {
		return 0, errors.New("token not active yet")
	}
	if c.ExpiresAt == nil || c.ExpiresAt.Before(now) {
		return 0, errors.New("token has expired")
	}
	return idFromString(c.Subject)
}
