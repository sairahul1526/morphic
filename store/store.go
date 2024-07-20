package store

import (
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/raystack/salt/db"
	"github.com/sairahul1526/morphic/metrics"
)

//go:embed migrations/*.sql
var migrationFs embed.FS

func Client(cfg db.Config) (*db.Client, error) {
	return db.New(cfg)
}

func Migrate(cfg db.Config) error {
	return db.RunMigrations(cfg, migrationFs, "migrations")
}

func Rollback(cfg db.Config) error {
	return db.RunRollback(cfg, migrationFs, "migrations")
}

func CheckPostgresError(err error) (error, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("%w [%s]", errDuplicateKey, pgErr.Detail), errDuplicateKey
		case pgerrcode.CheckViolation:
			return fmt.Errorf("%w [%s]", errCheckViolation, pgErr.Detail), errCheckViolation
		case pgerrcode.ForeignKeyViolation:
			return fmt.Errorf("%w [%s]", errForeignKeyViolation, pgErr.Detail), errForeignKeyViolation
		case pgerrcode.InvalidTextRepresentation:
			return fmt.Errorf("%w [%s]", errInvalidTexRepresentation, pgErr.Detail), errInvalidTexRepresentation
		}
	}
	return err, errors.New("unknown error")
}

func RecordMetrics(model, operation string, startTime time.Time, err error) {
	status := "success"
	errString := ""
	if err != nil {
		status = "error"
		_, e := CheckPostgresError(err)
		errString = e.Error()
	}

	metrics.DBRequests.WithLabelValues(model, operation, status, errString).Inc()
	metrics.DBLatency.WithLabelValues(model, operation, status, errString).Observe(float64(time.Since(startTime).Milliseconds()))
}
