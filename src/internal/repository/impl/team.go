package impl

import (
	"silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/model"

	"gorm.io/gorm"
)

type TeamRepository struct{}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}

func (teamRepo *TeamRepository) Create(conn abstract.IDBConnection, team *domain.Team) error {
	db := conn.Get()

	teamDAO := &model.Team{}
	teamDAO, err := teamDAO.ToModel(team)
	if err != nil {
		return err
	}

	return db.Save(teamDAO).Error
}

func (teamRepo *TeamRepository) GetByID(conn abstract.IDBConnection, id int) (*domain.Team, error) {
	db := conn.Get()

	var teamDAO model.Team
	err := db.Where("id= ?", id).First(&teamDAO).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return teamDAO.ToDomain()
}
