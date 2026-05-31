package abstract

import (
	"silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type ITeamRepository interface {
	Create(conn abstract.IDBConnection, team *domain.Team) error
	GetByID(conn abstract.IDBConnection, id int) (*domain.Team, error)
}
