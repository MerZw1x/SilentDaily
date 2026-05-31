package impl

import (
	"errors"

	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
	"silent/src/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AiApiRepository struct{}

func NewAiApiRepository() *AiApiRepository {
	return &AiApiRepository{}
}

func (r *AiApiRepository) InsertIfNotExist(conn db.IDBConnection, tokens []string) error {
	records := make([]model.AiApi, 0, len(tokens))
	for _, token := range tokens {
		hash := model.HashToken(token)
		dao, err := model.NewAiApiModel(hash)
		if err != nil {
			return err
		}
		records = append(records, *dao)
	}
	return conn.Get().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "hash"}},
			DoNothing: true,
		}).
		Create(&records).Error
}

func (r *AiApiRepository) IncreaseRequests(conn db.IDBConnection, token string) error {
	hash := model.HashToken(token)
	result := conn.Get().
		Model(&model.AiApi{}).
		Where("hash = ?", hash).
		Update("requests", gorm.Expr("requests + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ai_api token not found")
	}
	return nil
}

func (r *AiApiRepository) GetAllRequestsCount(conn db.IDBConnection, tokens []string) ([]*domain.AiApi, error) {
	hashMap := make(map[string]string, len(tokens))
	hashes := make([]string, len(tokens))
	for i, token := range tokens {
		h := model.HashToken(token)
		hashMap[h] = token
		hashes[i] = h
	}

	var daos []model.AiApi
	err := conn.Get().Where("hash IN ?", hashes).Find(&daos).Error
	if err != nil {
		return nil, err
	}

	result := make([]*domain.AiApi, len(daos))
	for i, dao := range daos {
		d, err := dao.ToDomain(hashMap[dao.Hash])
		if err != nil {
			return nil, err
		}
		result[i] = d
	}
	return result, nil
}
