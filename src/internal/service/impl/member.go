package impl

import (
	"errors"
	aconn "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/repository/abstract"
	"time"
)

type MemberService struct {
	Conn             aconn.IDBConnection
	MemberRepository abstract.IMemberRepository
}

func NewMemberRepository(conn aconn.IDBConnection, memberRepository abstract.IMemberRepository) *MemberService {
	return &MemberService{
		Conn:             conn,
		MemberRepository: memberRepository,
	}
}

func (service *MemberService) Register(telegramUserID int, name string, teamID int, isLead bool) error {
	member, err := service.MemberRepository.GetByTelegramID(service.Conn, telegramUserID)
	if err != nil {
		return err
	}
	if member != nil {
		return errors.New("the user already exists")
	}

	member = &domain.Member{
		TeamID:         teamID,
		TelegramUserID: telegramUserID,
		Name:           name,
		IsLead:         isLead,
		CreatedAt:      time.Now(),
	}

	err = service.MemberRepository.Create(service.Conn, member)
	if err != nil {
		return err
	}

	return nil
}

func (service *MemberService) GetByTelegramID(telegramUserID int) (*domain.Member, error) {
	member, err := service.MemberRepository.GetByTelegramID(service.Conn, telegramUserID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("the user does not exist")
	}

	return member, nil
}
