package employee

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/raystack/salt/db"
)

func MockClient(mockDB *sql.DB) *db.Client {
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return &db.Client{DB: sqlxDB}
}

// implement testcases
