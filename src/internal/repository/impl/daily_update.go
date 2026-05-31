package impl

import (
	"errors"
	"time"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/model"

	"gorm.io/gorm"
)

type DailyUpdateRepository struct{}

func NewDailyUpdateRepository() *DailyUpdateRepository {
	return &DailyUpdateRepository{}
}

func (r *DailyUpdateRepository) Create(conn db.IDBConnection, update *domain.DailyUpdate) error {
	dao := &model.DailyUpdate{}
	dao, err := dao.ToModel(update)
	if err != nil {
		return err
	}
	return conn.Get().Save(dao).Error
}

func (r *DailyUpdateRepository) GetOneQueued(conn db.IDBConnection) (*domain.DailyUpdate, error) {
	var dao model.DailyUpdate
	err := conn.Get().
		Raw(`SELECT * FROM silentdaily.daily_updates WHERE status = ? LIMIT 1 FOR UPDATE SKIP LOCKED`, "queued").
		Scan(&dao).Error
	if err != nil {
		return nil, err
	}
	if dao.ID == 0 {
		return nil, nil
	}
	return dao.ToDomain()
}

func (r *DailyUpdateRepository) SetStatus(conn db.IDBConnection, id int, status string) error {
	result := conn.Get().
		Model(&model.DailyUpdate{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("daily_update not found")
	}
	return nil
}

func (r *DailyUpdateRepository) SetStatusAndIncrementAttempts(conn db.IDBConnection, id int, status string) error {
	result := conn.Get().
		Model(&model.DailyUpdate{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":   status,
			"attempts": gorm.Expr("attempts + 1"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("daily_update not found")
	}
	return nil
}

func (r *DailyUpdateRepository) GetByTeamIDAndDate(conn db.IDBConnection, teamID int, date time.Time) ([]*domain.DailyUpdate, error) {
	var daos []model.DailyUpdate
	err := conn.Get().
		Where("team_id = ? AND DATE(created_at) = DATE(?)", teamID, date).
		Find(&daos).Error
	if err != nil {
		return nil, err
	}
	var m model.DailyUpdate
	return m.ToDomainSlice(daos)
}
