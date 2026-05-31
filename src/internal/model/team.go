package model

import (
	"silent/src/internal/domain"
	"time"
)

type Team struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (t *Team) ToDomain() (*domain.Team, error) {
	domain, err := ToDomain[Team, domain.Team](t)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (t *Team) ToModel(team *domain.Team) (*Team, error) {
	modelTeam, err := ToModel[Team, domain.Team](team)
	if err != nil {
		return nil, err
	}
	return modelTeam, nil
}
