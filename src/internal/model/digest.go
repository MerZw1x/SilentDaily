package model

import (
	"silent/src/internal/domain"
	"time"
)

type Digest struct {
	ID         int       `gorm:"column:id;primaryKey"`
	TeamID     int       `gorm:"column:team_id;not null"`
	Date       time.Time `gorm:"column:date;not null"`
	LeadDigest string    `gorm:"column:lead_digest;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (d *Digest) ToDomain() (*domain.Digest, error) {
	return ToDomain[Digest, domain.Digest](d)
}

func (d *Digest) ToModel(digest *domain.Digest) (*Digest, error) {
	return ToModel[Digest, domain.Digest](digest)
}
