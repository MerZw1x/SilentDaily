package impl

import (
	"errors"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/model"
)

type StructuredUpdateRepository struct{}

func NewStructuredUpdateRepository() *StructuredUpdateRepository {
	return &StructuredUpdateRepository{}
}

func (r *StructuredUpdateRepository) Create(conn db.IDBConnection, update *domain.StructuredUpdate) error {
	dao := &model.StructuredUpdate{}
	dao, err := dao.ToModel(update)
	if err != nil {
		return err
	}
	return conn.Get().Save(dao).Error
}

func (r *StructuredUpdateRepository) GetByDailyUpdateID(conn db.IDBConnection, dailyUpdateID int) (*domain.StructuredUpdate, error) {
	var dao model.StructuredUpdate
	err := conn.Get().
		Where("daily_update_id = ?", dailyUpdateID).
		First(&dao).Error
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dao.ToDomain()
}

func (r *StructuredUpdateRepository) GetByTeamIDAndDate(conn db.IDBConnection, teamID int, dailyUpdateIDs []int) ([]*domain.StructuredUpdate, error) {
	if len(dailyUpdateIDs) == 0 {
		return []*domain.StructuredUpdate{}, nil
	}
	var daos []model.StructuredUpdate
	err := conn.Get().
		Where("daily_update_id IN ?", dailyUpdateIDs).
		Find(&daos).Error
	if err != nil {
		return nil, err
	}
	var m model.StructuredUpdate
	return m.ToDomainSlice(daos)
}
