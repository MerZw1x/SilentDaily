package model

import (
	"silent/src/internal/domain"
	"time"
)

type Member struct {
	ID             int       `gorm:"column:id;primaryKey"`
	TeamID         int       `gorm:"column:team_id;not null"`
	TelegramUserID int       `gorm:"column:telegram_user_id;uniqueIndex;not null"`
	Name           string    `gorm:"column:name;not null"`
	IsLead         bool      `gorm:"column:is_lead;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (m *Member) ToDomain() (*domain.Member, error) {
	domain, err := ToDomain[Member, domain.Member](m)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (m *Member) ToModel(member *domain.Member) (*Member, error) {
	modelMember, err := ToModel[Member, domain.Member](member)
	if err != nil {
		return nil, err
	}
	return modelMember, nil
}
