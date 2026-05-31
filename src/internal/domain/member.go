package domain

import "time"

type Member struct {
	ID             int
	TeamID         int
	TelegramUserID int
	Name           string
	IsLead         bool
	CreatedAt      time.Time
}
