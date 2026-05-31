package impl

import (
	"errors"
	"time"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/model"

	"gorm.io/gorm/clause"
)

type DigestRepository struct{}

func NewDigestRepository() *DigestRepository {
	return &DigestRepository{}
}

func (r *DigestRepository) Upsert(conn db.IDBConnection, digest *domain.Digest) error {
	dao := &model.Digest{}
	dao, err := dao.ToModel(digest)
	if err != nil {
		return err
	}
	return conn.Get().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "team_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"lead_digest"}),
		}).
		Create(dao).Error
}

func (r *DigestRepository) GetByTeamIDAndDate(conn db.IDBConnection, teamID int, date time.Time) (*domain.Digest, error) {
	var dao model.Digest
	err := conn.Get().
		Where("team_id = ? AND date = ?", teamID, date.Truncate(24*time.Hour)).
		First(&dao).Error
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dao.ToDomain()
}
