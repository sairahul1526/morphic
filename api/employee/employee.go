package employee

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	api "github.com/sairahul1526/morphic/api"
	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/pkg/errors"
	"github.com/sairahul1526/morphic/pkg/validator"
	"go.uber.org/zap"
)

type EmployeeService interface {
	Create(ctx context.Context, employee entities.Employee) (entities.Employee, errors.Error)
	Delete(ctx context.Context, ids []string) errors.Error
	GetSummary(ctx context.Context, filters entities.SummaryFilter) (entities.Summary, errors.Error)
	GetSummaryByDepartment(ctx context.Context, filters entities.SummaryFilter) (map[string]entities.Summary, errors.Error)
	GetSummaryByDepartmentAndSubDepartment(ctx context.Context, filters entities.SummaryFilter) (map[string]map[string]entities.Summary, errors.Error)
}

type EmployeeHandlers struct {
	svc EmployeeService
}

func NewEmployeeHandlers(svc EmployeeService) *EmployeeHandlers {
	return &EmployeeHandlers{svc}
}

// @Tags Employee
// @Summary Add a new employee
// @Router /api/v1/employees [post]
// @Param body body Request true "Request Body"
// @Success 201 {object} ReadResponse
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @failure 400 {object} api.ErrorResponse
// @failure 404 {object} api.ErrorResponse
// @failure 403 {object} api.ErrorResponse
// @failure 409 {object} api.ErrorResponse
// @failure 500 {object} api.ErrorResponse
// @failure 501 {object} api.ErrorResponse
func (h *EmployeeHandlers) CreateEmployee(c *gin.Context) {
	var body Request
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("failed to bind request body", zap.Error(err))
		api.HarmonizeErrorResponse(c, errors.Error{Cause: errors.ErrCodeBadRequest, Message: err.Error(), Code: errors.ErrCodeBadRequest})
		return
	}

	validate := validator.GetValidator()
	if validationErr := validate.Struct(body); validationErr != nil {
		logger.Error("failed to create employee", []zap.Field{
			zap.Any("employee", body),
			zap.Error(validationErr),
		}...)
		api.HarmonizeErrorResponse(c, errors.Error{Cause: errors.ErrCodeBadRequest, Message: validationErr.Error(), Code: errors.ErrCodeInvalidRequest})
		return
	}

	employee, err := h.svc.Create(c, body.toDomain())
	if !err.IsEmpty() {
		logger.Error("failed to create employee", []zap.Field{
			zap.Any("employee", body),
			zap.Error(err),
		}...)
		api.HarmonizeErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, NewReadResponse(employee))
}

// @Tags Employee
// @Summary Delete a employee
// @Router /api/v1/employees [delete]
// @Param id query string true "Employee IDs"
// @Produce json
// @Success 200
// @Security ApiKeyAuth
// @failure 400 {object} api.ErrorResponse
// @failure 404 {object} api.ErrorResponse
// @failure 403 {object} api.ErrorResponse
// @failure 409 {object} api.ErrorResponse
// @failure 500 {object} api.ErrorResponse
// @failure 501 {object} api.ErrorResponse
func (h *EmployeeHandlers) DeleteEmployee(c *gin.Context) {
	err := h.svc.Delete(c, strings.Split(c.Query("id"), ","))
	if !err.IsEmpty() {
		logger.Error("failed to delete employee", []zap.Field{
			zap.String("id", c.Query("id")),
			zap.Error(err),
		}...)
		api.HarmonizeErrorResponse(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
