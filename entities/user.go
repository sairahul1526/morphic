package entities

import (
	"time"

	"github.com/sairahul1526/morphic/constant"
)

type User struct {
	ID        string
	Username  string
	Password  string
	Status    constant.UserStatus
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type UserFilter struct {
	Username string
	Password string
}
