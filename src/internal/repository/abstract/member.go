package abstract

import (
	"silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type IMemberRepository interface {
	Create(conn abstract.IDBConnection, member *domain.Member) error
	GetByTelegramID(conn abstract.IDBConnection, telegramUserID int64) (*domain.Member, error)
	GetByTeamID(conn abstract.IDBConnection, teamID int) ([]*domain.Member, error)
}
