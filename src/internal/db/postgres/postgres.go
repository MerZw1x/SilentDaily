package postgres

import (
	"fmt"

	"silent/src/internal/db/abstract"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PostgresDBConnection struct {
	conn *gorm.DB
}

func NewPostgresConnection(dsn string) *PostgresDBConnection {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "silentdaily.",
			SingularTable: true,
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Can not connect to postgres DB: %s", err.Error()))
	}

	return &PostgresDBConnection{conn: db}
}

func (c *PostgresDBConnection) Get() *gorm.DB {
	return c.conn
}

func (c *PostgresDBConnection) BeginTx() abstract.IDBConnection {
	return &PostgresDBConnection{conn: c.conn.Begin()}
}

func (c *PostgresDBConnection) Commit() error {
	return c.conn.Commit().Error
}

func (c *PostgresDBConnection) Rollback() {
	c.conn.Rollback()
}
