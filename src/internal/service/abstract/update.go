package abstract

import (
	"silent/src/internal/db/abstract"
)

type IUpdateService interface {
	Submit(conn abstract.IDBConnection, telegramUserID int, rawText string) error
}
