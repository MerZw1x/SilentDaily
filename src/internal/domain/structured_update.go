package domain

import "time"

type StructuredUpdate struct {
	ID            int
	DailyUpdateID int
	Progress      []string
	Plans         []string
	Blockers      []string
	CreatedAt     time.Time
}
