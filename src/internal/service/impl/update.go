package impl

import (
	"errors"
	aconn "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/repository/abstract"
)

type UpdateService struct {
	Conn                  aconn.IDBConnection
	MemberRepository      abstract.IMemberRepository
	DailyUpdateRepository abstract.IDailyUpdateRepository
}

func NewUpdateService(conn aconn.IDBConnection, memberRepository abstract.IMemberRepository, dailyUpdateRepository abstract.IDailyUpdateRepository) *UpdateService {
	return &UpdateService{
		Conn:                  conn,
		MemberRepository:      memberRepository,
		DailyUpdateRepository: dailyUpdateRepository,
	}
}

func (service *UpdateService) Submit(telegramUserID int, rawText string) error {
	member, err := service.MemberRepository.GetByTelegramID(service.Conn, telegramUserID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("member not registered")
	}

	update := createDailyUpdate(member.ID, member.TeamID, rawText)
	err = service.DailyUpdateRepository.Create(service.Conn, update)
	if err != nil {
		return err
	}
	return nil
}

func createDailyUpdate(memberID int, teamID int, rawText string) *domain.DailyUpdate {
	return &domain.DailyUpdate{
		MemberID: memberID,
		TeamID:   teamID,
		RawText:  rawText,
		Status:   "queued",
	}
}
