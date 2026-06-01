package impl

import (
	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	repoabstract "silent/src/internal/repository/abstract"
)

type TeamService struct {
	conn     db.IDBConnection
	teamRepo repoabstract.ITeamRepository
}

func NewTeamService(conn db.IDBConnection, teamRepo repoabstract.ITeamRepository) *TeamService {
	return &TeamService{conn: conn, teamRepo: teamRepo}
}

func (s *TeamService) Create(name string) (*domain.Team, error) {
	team := &domain.Team{Name: name}
	if err := s.teamRepo.Create(s.conn, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) GetByID(id int) (*domain.Team, error) {
	return s.teamRepo.GetByID(s.conn, id)
}
