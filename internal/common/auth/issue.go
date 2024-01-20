package auth

import "time"

type Issuer interface {
	IssueToken(subjectID int64, notBefore, expiresAt time.Time) (string, error)
}

type staticIssuer struct {
	privateKey any
}

func NewStaticIssuer(pk any) Issuer {
	return &staticIssuer{privateKey: pk}
}

func (i *staticIssuer) IssueToken(subjectID int64, notBefore, expiresAt time.Time) (string, error) {
	return IssueToken(subjectID, notBefore, expiresAt, i.privateKey)
}
