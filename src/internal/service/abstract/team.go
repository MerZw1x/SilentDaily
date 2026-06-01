package abstract

import "silent/src/internal/domain"

type ITeamService interface {
	Create(name string) (*domain.Team, error)
	GetByID(id int) (*domain.Team, error)
}
