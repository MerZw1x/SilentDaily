package impl

import (
	"silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/model"

	"gorm.io/gorm"
)

type MemberRepository struct{}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{}
}

func (memberRepo *MemberRepository) Create(conn abstract.IDBConnection, member *domain.Member) error {
	db := conn.Get()

	memberDAO := &model.Member{}
	memberDAO, err := memberDAO.ToModel(member)
	if err != nil {
		return err
	}

	return db.Save(memberDAO).Error
}

func (memberRepo *MemberRepository) GetByTelegramID(conn abstract.IDBConnection, telegramUserID int) (*domain.Member, error) {
	db := conn.Get()

	var memberDAO model.Member
	err := db.Where("telegram_user_id = ?", telegramUserID).First(&memberDAO).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return memberDAO.ToDomain()
}

func (memberRepo *MemberRepository) GetByTeamID(conn abstract.IDBConnection, teamID int) ([]*domain.Member, error) {
	db := conn.Get()

	var memberDAOs []model.Member
	err := db.Find(&memberDAOs).Error

	if err != nil {
		return nil, err
	}

	var modelObj model.Member
	return modelObj.ToDomainSlice(memberDAOs)
}
