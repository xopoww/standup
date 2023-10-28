package auth

import "time"

type Issuer interface {
	IssueToken(subjectId string, notBefore, expiresAt time.Time) (string, error)
}

type staticIssuer struct {
	privateKey any
}

func NewStaticIssuer(pk any) Issuer {
	return &staticIssuer{privateKey: pk}
}

func (i *staticIssuer) IssueToken(subjectId string, notBefore, expiresAt time.Time) (string, error) {
	return IssueToken(subjectId, notBefore, expiresAt, i.privateKey)
}
