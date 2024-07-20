package user

import (
	"time"

	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
)

type ReadResponse struct {
	ID        string              `json:"id" example:"47435e2b-d8c4-41ff-9de9-2be3bfc92276"`
	Username  string              `json:"username" example:"john.doe"`
	Status    constant.UserStatus `json:"status" example:"Active"`
	CreatedAt time.Time           `json:"created_at" example:"2024-05-27T06:50:20.056549Z"`
	CreatedBy string              `json:"created_by" example:"e31ab6f8-d359-4c6a-83c6-bfa32229bb01"`
	UpdatedAt time.Time           `json:"updated_at" example:"2024-05-27T06:50:20.056549Z"`
	UpdatedBy string              `json:"updated_by" example:"e31ab6f8-d359-4c6a-83c6-bfa32229bb01"`
}

func NewReadResponse(user entities.User) ReadResponse {
	return ReadResponse{
		ID:        user.ID,
		Username:  user.Username,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}
}
