package employee

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/sairahul1526/morphic/api"
	"github.com/sairahul1526/morphic/logger"
	"go.uber.org/zap"
)

// @Tags Summary
// @Summary Get salary summary
// @Router /api/v1/employees/summary [get]
// @Param on_contract query string false "Filter by contract type"
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Summary
// @failure 400 {object} api.ErrorResponse
// @failure 404 {object} api.ErrorResponse
// @failure 403 {object} api.ErrorResponse
// @failure 409 {object} api.ErrorResponse
// @failure 500 {object} api.ErrorResponse
// @failure 501 {object} api.ErrorResponse
func (h *EmployeeHandlers) GetSummary(c *gin.Context) {
	summary, err := h.svc.GetSummary(c, parseFilters(c))
	if !err.IsEmpty() {
		logger.Error("failed to get salary summary", []zap.Field{
			zap.Any("filters", parseFilters(c)),
			zap.Error(err),
		}...)
		api.HarmonizeErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, NewSummary(summary))
}

// @Tags Summary
// @Summary Get salary summary by department
// @Router /api/v1/employees/summary/department [get]
// @Param on_contract query string false "Filter by contract type"
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} DepartmentSummaryResponse
// @failure 400 {object} api.ErrorResponse
// @failure 404 {object} api.ErrorResponse
// @failure 403 {object} api.ErrorResponse
// @failure 409 {object} api.ErrorResponse
// @failure 500 {object} api.ErrorResponse
// @failure 501 {object} api.ErrorResponse
func (h *EmployeeHandlers) GetSummaryByDepartment(c *gin.Context) {
	result, err := h.svc.GetSummaryByDepartment(c, parseFilters(c))
	if !err.IsEmpty() {
		logger.Error("failed to get summary by department", []zap.Field{
			zap.Any("filters", parseFilters(c)),
			zap.Error(err),
		}...)
		api.HarmonizeErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, NewSummaryByDepartment(result))
}

// @Tags Summary
// @Summary Get salary summary by department and sub-department
// @Router /api/v1/employees/summary/subdepartment [get]
// @Param on_contract query string false "Filter by contract type"
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} DepartmentSummaryResponse
// @failure 400 {object} api.ErrorResponse
// @failure 404 {object} api.ErrorResponse
// @failure 403 {object} api.ErrorResponse
// @failure 409 {object} api.ErrorResponse
// @failure 500 {object} api.ErrorResponse
// @failure 501 {object} api.ErrorResponse
func (h *EmployeeHandlers) GetSummaryByDepartmentAndSubDepartment(c *gin.Context) {
	result, err := h.svc.GetSummaryByDepartmentAndSubDepartment(c, parseFilters(c))
	if !err.IsEmpty() {
		logger.Error("failed to get salary summary by department and sub-department", []zap.Field{
			zap.Any("filters", parseFilters(c)),
			zap.Error(err),
		}...)
		api.HarmonizeErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, NewSummaryByDepartmentAndSubDepartment(result))
}
