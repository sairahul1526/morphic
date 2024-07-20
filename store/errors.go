package store

import (
	"github.com/sairahul1526/morphic/pkg/errors"
)

var (
	errDuplicateKey             = errors.New("duplicate key")
	errCheckViolation           = errors.New("check constraint violation")
	errForeignKeyViolation      = errors.New("foreign key violation")
	errInvalidTexRepresentation = errors.New("invalid input syntax type")
)
