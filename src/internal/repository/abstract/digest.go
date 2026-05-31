package abstract

import (
	"time"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type IDigestRepository interface {
	Upsert(conn db.IDBConnection, digest *domain.Digest) error
	GetByTeamIDAndDate(conn db.IDBConnection, teamID int, date time.Time) (*domain.Digest, error)
}
