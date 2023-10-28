package models

import "time"

type Message struct {
	ID        string
	OwnerID   string
	Text      string
	CreatedAt time.Time
}
