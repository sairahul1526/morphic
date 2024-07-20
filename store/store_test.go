package store

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/raystack/salt/db"
	"github.com/stretchr/testify/assert"
)

func TestClientEmptyDriver(t *testing.T) {
	_, err := Client(db.Config{
		Driver: "",
		URL:    "postgres://postgres:password@localhost:5432/morphic?sslmode=disable",
	})
	assert.EqualError(t, err, `sql: unknown driver "" (forgotten import?)`)
}

func TestMigrateEmptyURL(t *testing.T) {
	err := Migrate(db.Config{
		Driver: "postgres",
		URL:    "",
	})
	assert.EqualError(t, err, `URL cannot be empty`)
}

func TestRollbackEmptyURL(t *testing.T) {
	err := Rollback(db.Config{
		Driver: "postgres",
		URL:    "",
	})
	assert.EqualError(t, err, `URL cannot be empty`)
}

func TestCheckPostgresErrorUniqueViolation(t *testing.T) {
	result, _ := CheckPostgresError(&pgconn.PgError{Code: pgerrcode.UniqueViolation, Detail: "Duplicate key value violates unique constraint"})
	assert.EqualError(t, result, fmt.Errorf("%w [%s]", errDuplicateKey, "Duplicate key value violates unique constraint").Error())
	assert.Equal(t, "duplicate key [Duplicate key value violates unique constraint]", result.Error())
}

func TestCheckPostgresErrorCheckViolation(t *testing.T) {
	result, _ := CheckPostgresError(&pgconn.PgError{Code: pgerrcode.CheckViolation, Detail: "New row violates check constraint"})
	assert.EqualError(t, result, fmt.Errorf("%w [%s]", errCheckViolation, "New row violates check constraint").Error())
	assert.Equal(t, "check constraint violation [New row violates check constraint]", result.Error())
}

func TestCheckPostgresErrorForeignKeyViolation(t *testing.T) {
	result, _ := CheckPostgresError(&pgconn.PgError{Code: pgerrcode.ForeignKeyViolation, Detail: "Key value violates foreign key constraint"})
	assert.EqualError(t, result, fmt.Errorf("%w [%s]", errForeignKeyViolation, "Key value violates foreign key constraint").Error())
	assert.Equal(t, "foreign key violation [Key value violates foreign key constraint]", result.Error())
}

func TestCheckPostgresErrorInvalidTextRepresentation(t *testing.T) {
	result, _ := CheckPostgresError(&pgconn.PgError{Code: pgerrcode.InvalidTextRepresentation, Detail: "Invalid input syntax for type"})
	assert.EqualError(t, result, fmt.Errorf("%w [%s]", errInvalidTexRepresentation, "Invalid input syntax for type").Error())
	assert.Equal(t, "invalid input syntax type [Invalid input syntax for type]", result.Error())
}

func TestCheckPostgresErrorOtherError(t *testing.T) {
	result, _ := CheckPostgresError(errors.New("some other error"))
	assert.EqualError(t, result, "some other error")
}
