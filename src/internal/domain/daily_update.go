package domain

import "time"

type DailyUpdate struct {
	ID        int
	MemberID  int
	TeamID    int
	RawText   string
	Status    string
	Attempts  int
	CreatedAt time.Time
}
