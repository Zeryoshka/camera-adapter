package confapi

import (
	"database/sql"
)

type ManagedDevice struct {
	Id       int
	Name     string
	Host     string
	Login    string
	Password string
}

// DB is a database of stock trades.
type DB struct {
	sql    *sql.DB
	stmt   *sql.Stmt
	buffer []ManagedDevice
}
