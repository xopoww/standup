package models

import "time"

type Message struct {
	ID        string
	OwnerID   int64
	Text      string
	CreatedAt time.Time
}
