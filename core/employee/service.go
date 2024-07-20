package employee

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s Service) Create(ctx context.Context, employee entities.Employee) (result entities.Employee, err errors.Error) {

	tx := ctx.Value(constant.DBTrxKey).(*sqlx.Tx)
	defer func() {
		if tx.Tx != nil && !err.IsEmpty() {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Error("failed to rollback transaction", []zap.Field{
					zap.Error(rollbackErr),
				}...)
			}
		}
	}()

	employee.ID = uuid.NewString()

	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		logger.Error("failed to get gin context")
	} else {
		employee.CreatedBy = ginCtx.GetHeader("user_id")
		employee.UpdatedBy = ginCtx.GetHeader("user_id")
	}

	// add to transaction
	createdEmployee, dbErr := s.repo.Create(ctx, tx, employee)
	if dbErr != nil {
		return entities.Employee{}, errors.Error{Cause: errors.ErrCodeInternalServer, Message: dbErr.Error(), Code: errors.ErrDBCreate}
	}

	if err := tx.Commit(); err != nil {
		return entities.Employee{}, errors.Error{Cause: errors.ErrCodeInternalServer, Message: err.Error(), Code: errors.ErrDBCommit}
	}

	return createdEmployee, errors.Error{}
}

func (s Service) Delete(ctx context.Context, ids []string) errors.Error {
	err := s.repo.Delete(ctx, ids)
	if err != nil {
		return errors.Error{Cause: errors.ErrCodeInternalServer, Message: err.Error(), Code: errors.ErrDBDelete}
	}
	return errors.Error{}
}

func (s Service) GetSummary(ctx context.Context, filters entities.SummaryFilter) (entities.Summary, errors.Error) {
	summary, err := s.repo.GetSummary(ctx, filters)
	if err != nil {
		return entities.Summary{}, errors.Error{Cause: errors.ErrCodeInternalServer, Message: err.Error(), Code: errors.ErrCodeUnknown}
	}
	return summary, errors.Error{}
}

func (s Service) GetSummaryByDepartment(ctx context.Context, filters entities.SummaryFilter) (map[string]entities.Summary, errors.Error) {
	summaries, err := s.repo.GetSummaryByDepartment(ctx, filters)
	if err != nil {
		return map[string]entities.Summary{}, errors.Error{Cause: errors.ErrCodeInternalServer, Message: err.Error(), Code: errors.ErrCodeUnknown}
	}
	return summaries, errors.Error{}
}

func (s Service) GetSummaryByDepartmentAndSubDepartment(ctx context.Context, filters entities.SummaryFilter) (map[string]map[string]entities.Summary, errors.Error) {
	summaries, err := s.repo.GetSummaryByDepartmentAndSubDepartment(ctx, filters)
	if err != nil {
		return map[string]map[string]entities.Summary{}, errors.Error{Cause: errors.ErrCodeInternalServer, Message: err.Error(), Code: errors.ErrCodeUnknown}
	}
	return summaries, errors.Error{}
}
