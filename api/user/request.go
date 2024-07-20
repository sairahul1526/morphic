package user

import (
	"github.com/sairahul1526/morphic/entities"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"john.doe"`
	Password string `json:"password" validate:"required" example:"password"`
}

func (r LoginRequest) toDomain() entities.UserFilter {
	return entities.UserFilter{
		Username: r.Username,
		Password: r.Password,
	}
}
