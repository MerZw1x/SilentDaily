package abstract

import (
	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type IStructuredUpdateRepository interface {
	Create(conn db.IDBConnection, update *domain.StructuredUpdate) error
	GetByDailyUpdateID(conn db.IDBConnection, dailyUpdateID int) (*domain.StructuredUpdate, error)
	GetByTeamIDAndDate(conn db.IDBConnection, teamID int, dailyUpdateIDs []int) ([]*domain.StructuredUpdate, error)
}
