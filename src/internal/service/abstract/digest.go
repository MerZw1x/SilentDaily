package abstract

import (
	"time"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type IDigestService interface {
	GetByTeamAndDate(conn db.IDBConnection, teamID int, date time.Time) (*domain.Digest, error)
}
