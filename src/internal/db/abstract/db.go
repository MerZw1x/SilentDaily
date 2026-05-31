package abstract

import "gorm.io/gorm"

type IDBConnection interface {
	Get() *gorm.DB
	BeginTx() IDBConnection
	Commit() error
	Rollback()
}
