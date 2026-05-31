package model

import (
	"silent/src/internal/domain"
	"time"

	"github.com/lib/pq"
)

type StructuredUpdate struct {
	ID            int            `gorm:"column:id;primaryKey"`
	DailyUpdateID int            `gorm:"column:daily_update_id;not null"`
	Progress      pq.StringArray `gorm:"column:progress;type:text[]"`
	Plans         pq.StringArray `gorm:"column:plans;type:text[]"`
	Blockers      pq.StringArray `gorm:"column:blockers;type:text[]"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
}

func (s *StructuredUpdate) ToDomain() (*domain.StructuredUpdate, error) {
	result := &domain.StructuredUpdate{
		ID:            s.ID,
		DailyUpdateID: s.DailyUpdateID,
		Progress:      []string(s.Progress),
		Plans:         []string(s.Plans),
		Blockers:      []string(s.Blockers),
		CreatedAt:     s.CreatedAt,
	}
	return result, nil
}

func (s *StructuredUpdate) ToModel(su *domain.StructuredUpdate) (*StructuredUpdate, error) {
	return &StructuredUpdate{
		ID:            su.ID,
		DailyUpdateID: su.DailyUpdateID,
		Progress:      pq.StringArray(su.Progress),
		Plans:         pq.StringArray(su.Plans),
		Blockers:      pq.StringArray(su.Blockers),
		CreatedAt:     su.CreatedAt,
	}, nil
}

func (s *StructuredUpdate) ToDomainSlice(updates []StructuredUpdate) ([]*domain.StructuredUpdate, error) {
	result := make([]*domain.StructuredUpdate, len(updates))
	for i := range updates {
		d, err := updates[i].ToDomain()
		if err != nil {
			return nil, err
		}
		result[i] = d
	}
	return result, nil
}
