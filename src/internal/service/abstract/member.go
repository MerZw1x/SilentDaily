package abstract

import (
	"silent/src/internal/domain"
)

type IMemberService interface {
	Register(telegramUserID int, name string, teamID int, isLead bool) error
	GetByTelegramID(telegramUserID int) (*domain.Member, error)
}
