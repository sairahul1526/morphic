package user

import (
	"context"

	"github.com/sairahul1526/morphic/entities"
)

type Repository interface {
	Get(context.Context, entities.UserFilter) (entities.User, error)
}
