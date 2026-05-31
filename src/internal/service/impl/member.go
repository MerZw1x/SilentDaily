package impl

import (
	"errors"
	aconn "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/repository/abstract"
	"time"
)

type MemberService struct {
	MemberRepository abstract.IMemberRepository
}

func NewMemberRepository(memberRepository abstract.IMemberRepository) *MemberService {
	return &MemberService{
		MemberRepository: memberRepository,
	}
}

func (service *MemberService) Register(conn aconn.IDBConnection, telegramUserID int, name string, teamID int, isLead bool) error {
	member, err := service.MemberRepository.GetByTelegramID(conn, telegramUserID)
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

	err = service.MemberRepository.Create(conn, member)
	if err != nil {
		return err
	}

	return nil
}

func (service *MemberService) GetByTelegramID(conn aconn.IDBConnection, telegramUserID int) (*domain.Member, error) {
	member, err := service.MemberRepository.GetByTelegramID(conn, telegramUserID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("the user does not exist")
	}

	return member, nil
}
