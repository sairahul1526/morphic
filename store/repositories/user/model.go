package user

import (
	"time"

	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
)

type UserModel struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}

func (u *UserModel) ToDomain() entities.User {
	return entities.User{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Status:    constant.UserStatus(u.Status),
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: u.UpdatedBy,
	}
}
