package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/raystack/salt/db"
	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/store"
	"go.uber.org/zap"
)

type UserRepository struct {
	db *db.Client
}

func NewRepository(db *db.Client) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Get(ctx context.Context, filters entities.UserFilter) (user entities.User, err error) {
	startTime := time.Now()
	defer store.RecordMetrics("user", "get", startTime, err)

	userFilter := UserFilter{
		Username: filters.Username,
		Password: filters.Password,
	}
	filterQuery, args := userFilter.ToQuery()

	query := `SELECT * FROM users ` + filterQuery

	var userModel UserModel
	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&userModel)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, nil
		}
		logger.Error("failed to get user", []zap.Field{
			zap.Any("filters", filters),
			zap.Error(err),
		}...)
		return
	}

	return userModel.ToDomain(), nil
}
