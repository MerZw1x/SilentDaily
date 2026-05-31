package abstract

import (
	"time"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type IDailyUpdateRepository interface {
	Create(conn db.IDBConnection, update *domain.DailyUpdate) error
	GetOneQueued(conn db.IDBConnection) (*domain.DailyUpdate, error)
	SetStatus(conn db.IDBConnection, id int, status string) error
	SetStatusAndIncrementAttempts(conn db.IDBConnection, id int, status string) error
	GetByTeamIDAndDate(conn db.IDBConnection, teamID int, date time.Time) ([]*domain.DailyUpdate, error)
}
