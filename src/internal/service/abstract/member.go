package abstract

import (
	"silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type IMemberService interface {
	Register(conn abstract.IDBConnection, telegramUserID int, name string, teamID int, isLead bool) error
	GetByTelegramID(conn abstract.IDBConnection, telegramUserID int) (*domain.Member, error)
}
