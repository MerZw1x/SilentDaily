package domain

import "time"

type Team struct {
	ID        int
	Name      string
	CreatedAt time.Time
}
