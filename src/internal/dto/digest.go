package dto

import "time"

type DigestResponse struct {
	TeamID     int       `json:"team_id"`
	Date       time.Time `json:"date"`
	LeadDigest string    `json:"lead_digest"`
}
