package employee_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sairahul1526/morphic/api/employee"
	"github.com/sairahul1526/morphic/api/employee/mocks"
	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmployeeHandlers_GetSummary(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.GET("/api/v1/employees/summary", handler.GetSummary)

	summary := entities.Summary{
		Mean: 50000,
		Min:  30000,
		Max:  70000,
	}

	mockSvc.On("GetSummary", mock.Anything, mock.AnythingOfType("entities.SummaryFilter")).Return(summary, errors.Error{})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), `{"mean":50000,"max":70000,"min":30000}`)

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_GetSummary_ServiceError(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.GET("/api/v1/employees/summary", handler.GetSummary)

	expectedError := errors.Error{Cause: errors.ErrCodeInternalServer, Message: "failed to get salary summary", Code: errors.ErrCodeInternalServer}
	mockSvc.On("GetSummary", mock.Anything, mock.AnythingOfType("entities.SummaryFilter")).Return(entities.Summary{}, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get salary summary")

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_GetSummaryByDepartment(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.GET("/api/v1/employees/summary/department", handler.GetSummaryByDepartment)

	departmentSummary := map[string]entities.Summary{
		"Engineering": {Mean: 60000, Min: 40000, Max: 80000},
	}

	mockSvc.On("GetSummaryByDepartment", mock.Anything, mock.AnythingOfType("entities.SummaryFilter")).Return(departmentSummary, errors.Error{})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/summary/department", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), `{"summaries":[{"department":"Engineering","summary":{"mean":60000,"max":80000,"min":40000}}]}`)

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_GetSummaryByDepartment_ServiceError(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.GET("/api/v1/employees/summary/department", handler.GetSummaryByDepartment)

	expectedError := errors.Error{Cause: errors.ErrCodeInternalServer, Message: "failed to get summary by department", Code: errors.ErrCodeInternalServer}
	mockSvc.On("GetSummaryByDepartment", mock.Anything, mock.AnythingOfType("entities.SummaryFilter")).Return(nil, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/summary/department", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `failed to get summary by department`)

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_GetSummaryByDepartmentAndSubDepartment(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.GET("/api/v1/employees/summary/subdepartment", handler.GetSummaryByDepartmentAndSubDepartment)

	subDepartmentSummary := map[string]map[string]entities.Summary{
		"Engineering": {
			"Platform": {Mean: 70000, Min: 50000, Max: 90000},
		},
	}

	mockSvc.On("GetSummaryByDepartmentAndSubDepartment", mock.Anything, mock.AnythingOfType("entities.SummaryFilter")).Return(subDepartmentSummary, errors.Error{})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/summary/subdepartment", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), `{"summaries":[{"department":"Engineering","summaries":[{"sub_department":"Platform","summary":{"mean":70000,"max":90000,"min":50000}}]}]}`)

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_GetSummaryByDepartmentAndSubDepartment_ServiceError(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.GET("/api/v1/employees/summary/subdepartment", handler.GetSummaryByDepartmentAndSubDepartment)

	expectedError := errors.Error{Cause: errors.ErrCodeInternalServer, Message: "failed to get summary by department and sub-department", Code: errors.ErrCodeInternalServer}
	mockSvc.On("GetSummaryByDepartmentAndSubDepartment", mock.Anything, mock.AnythingOfType("entities.SummaryFilter")).Return(nil, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/summary/subdepartment", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `failed to get summary by department and sub-department`)

	mockSvc.AssertExpectations(t)
}
