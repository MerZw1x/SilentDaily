package model

import (
	"silent/src/internal/domain"
	"time"
)

type DailyUpdate struct {
	ID        int       `gorm:"column:id;primaryKey"`
	MemberID  int       `gorm:"column:member_id;not null"`
	TeamID    int       `gorm:"column:team_id;not null"`
	RawText   string    `gorm:"column:raw_text;not null"`
	Status    string    `gorm:"column:status;not null;default:queued"`
	Attempts  int       `gorm:"column:attempts;not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (d *DailyUpdate) ToDomain() (*domain.DailyUpdate, error) {
	return ToDomain[DailyUpdate, domain.DailyUpdate](d)
}

func (d *DailyUpdate) ToModel(update *domain.DailyUpdate) (*DailyUpdate, error) {
	return ToModel[DailyUpdate, domain.DailyUpdate](update)
}

func (d *DailyUpdate) ToDomainSlice(updates []DailyUpdate) ([]*domain.DailyUpdate, error) {
	return ToDomainSlice[DailyUpdate, domain.DailyUpdate](updates)
}
