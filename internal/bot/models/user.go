package models

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	ID       int64
	Username string
}

func (u User) String() string {
	return fmt.Sprintf("%q (%d)", u.Username, u.ID)
}

func FromTG(tu *tgbotapi.User) *User {
	return &User{
		ID:       tu.ID,
		Username: tu.UserName,
	}
}
