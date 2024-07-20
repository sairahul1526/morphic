package user

import (
	"context"

	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/pkg/errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s Service) Get(ctx context.Context, filters entities.UserFilter) (result entities.User, err errors.Error) {
	user, dbErr := s.repo.Get(ctx, filters)
	if dbErr != nil {
		return entities.User{}, errors.Error{Cause: errors.ErrCodeInternalServer, Message: dbErr.Error(), Code: errors.ErrDBGet}
	}
	return user, errors.Error{}
}
