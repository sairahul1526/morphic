package employee_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sairahul1526/morphic/api/employee"
	"github.com/sairahul1526/morphic/api/employee/mocks"
	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestEmployeeHandlers_CreateEmployee(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.POST("/api/v1/employees", handler.CreateEmployee)

	employeeToCreate := entities.Employee{
		Name:          "John Doe",
		Currency:      constant.EmployeeCurrencyUSD,
		Salary:        100000,
		Department:    constant.EmployeeDepartmentEngineering,
		SubDepartment: constant.EmployeeSubDepartmentPlatform,
		OnContract:    true,
		Status:        constant.EmployeeStatusActive,
	}
	employeeToReturn := entities.Employee{
		ID:            "employee123",
		Name:          "John Doe",
		Currency:      constant.EmployeeCurrencyUSD,
		Salary:        100000,
		Department:    constant.EmployeeDepartmentEngineering,
		SubDepartment: constant.EmployeeSubDepartmentPlatform,
		OnContract:    true,
		Status:        constant.EmployeeStatusActive,
		CreatedBy:     "creator123",
		UpdatedBy:     "updater123",
	}

	mockSvc.On("Create", mock.Anything, employeeToCreate).Return(employeeToReturn, errors.Error{})

	body, _ := json.Marshal(employee.Request{
		Name:          "John Doe",
		Currency:      constant.EmployeeCurrencyUSD,
		Salary:        100000,
		Department:    constant.EmployeeDepartmentEngineering,
		SubDepartment: constant.EmployeeSubDepartmentPlatform,
		OnContract:    true,
		Status:        constant.EmployeeStatusActive,
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	fmt.Println(w.Body.String())
	assert.JSONEq(t, `{"id":"employee123","name":"John Doe","currency":"`+constant.EmployeeCurrencyUSD.String()+`","salary":100000,"department":"`+constant.EmployeeDepartmentEngineering.String()+`","sub_department":"`+constant.EmployeeSubDepartmentPlatform.String()+`","on_contract":true,"status":"`+constant.EmployeeStatusActive.String()+`","created_at":"0001-01-01T00:00:00Z","created_by":"creator123","updated_at":"0001-01-01T00:00:00Z","updated_by":"updater123"}`, w.Body.String())

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_CreateEmployee_BindJSONError(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.POST("/api/v1/employees", handler.CreateEmployee)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString("{invalid json}"))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `invalid character 'i' looking for beginning of object key string`)

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_CreateEmployee_ValidationError(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.POST("/api/v1/employees", handler.CreateEmployee)

	body, _ := json.Marshal(employee.Request{
		Name:          "",
		Currency:      "",
		Salary:        0,
		Department:    "",
		SubDepartment: "",
		OnContract:    false,
		Status:        "",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_CreateEmployee_ServiceError(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.POST("/api/v1/employees", handler.CreateEmployee)

	expectedError := errors.Error{Cause: errors.ErrCodeInternalServer, Message: "failed to create employee", Code: errors.ErrCodeInternalServer}
	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("entities.Employee")).Return(entities.Employee{}, expectedError)

	body, _ := json.Marshal(employee.Request{
		Name:          "John Doe",
		Currency:      constant.EmployeeCurrencyUSD,
		Salary:        100000,
		Department:    constant.EmployeeDepartmentEngineering,
		SubDepartment: constant.EmployeeSubDepartmentPlatform,
		OnContract:    true,
		Status:        constant.EmployeeStatusActive,
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":{"code":"`+errors.ErrCodeInternalServer+`","message":"failed to create employee"}}`, w.Body.String())

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_DeleteEmployee(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.DELETE("/api/v1/employees", handler.DeleteEmployee)

	mockSvc.On("Delete", mock.Anything, []string{"id1"}).Return(errors.Error{})

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/employees?id=id1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockSvc.AssertExpectations(t)
}

func TestEmployeeHandlers_DeleteEmployee_ServiceError(t *testing.T) {
	mockSvc := mocks.NewEmployeeService(t)
	handler := employee.NewEmployeeHandlers(mockSvc)

	r := setupRouter()
	r.DELETE("/api/v1/employees", handler.DeleteEmployee)

	expectedError := errors.Error{Cause: errors.ErrCodeInternalServer, Message: "failed to delete employee", Code: errors.ErrCodeInternalServer}
	mockSvc.On("Delete", mock.Anything, []string{"id1"}).Return(expectedError)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/employees?id=id1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":{"code":"`+errors.ErrCodeInternalServer+`","message":"failed to delete employee"}}`, w.Body.String())

	mockSvc.AssertExpectations(t)
}
