package abstract

import (
	db "silent/src/internal/db/abstract"
	"silent/src/internal/domain"
)

type IAiApiRepository interface {
	InsertIfNotExist(conn db.IDBConnection, tokens []string) error
	IncreaseRequests(conn db.IDBConnection, token string) error
	GetAllRequestsCount(conn db.IDBConnection, tokens []string) ([]*domain.AiApi, error)
}
