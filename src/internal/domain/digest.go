package domain

import "time"

type Digest struct {
	ID         int
	TeamID     int
	Date       time.Time
	LeadDigest string
	CreatedAt  time.Time
}
