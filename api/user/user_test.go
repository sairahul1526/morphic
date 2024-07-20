package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sairahul1526/morphic/api/user"
	"github.com/sairahul1526/morphic/api/user/mocks"
	"github.com/sairahul1526/morphic/config"
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

func TestUserHandlers_Login(t *testing.T) {
	mockSvc := mocks.NewUserService(t)
	cfg := config.Config{
		Auth: config.AuthConfig{
			Secret: "testsecret",
			Expiry: 3600,
		},
	}
	handler := user.NewUserHandlers(mockSvc, cfg)

	r := setupRouter()
	r.POST("/api/v1/users/login", handler.Login)

	userToReturn := entities.User{
		ID:        "user123",
		Username:  "testuser",
		Status:    constant.UserStatusActive,
		CreatedBy: "creator123",
		UpdatedBy: "updater123",
	}

	mockSvc.On("Get", mock.Anything, entities.UserFilter{
		Username: "testuser",
		Password: "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
	}).Return(userToReturn, errors.Error{})

	body, _ := json.Marshal(user.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("access_token"))
	assert.Equal(t, `{"id":"user123","username":"testuser","status":"Active","created_at":"0001-01-01T00:00:00Z","created_by":"creator123","updated_at":"0001-01-01T00:00:00Z","updated_by":"updater123"}`, w.Body.String())
	mockSvc.AssertExpectations(t)
}

func TestUserHandlers_Login_BindJSONError(t *testing.T) {
	mockSvc := mocks.NewUserService(t)
	cfg := config.Config{}
	handler := user.NewUserHandlers(mockSvc, cfg)

	r := setupRouter()
	r.POST("/api/v1/users/login", handler.Login)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBufferString("{invalid json}"))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `invalid character 'i' looking for beginning of object key string`)

	mockSvc.AssertExpectations(t)
}

func TestUserHandlers_Login_ValidationError(t *testing.T) {
	mockSvc := mocks.NewUserService(t)
	cfg := config.Config{}
	handler := user.NewUserHandlers(mockSvc, cfg)

	r := setupRouter()
	r.POST("/api/v1/users/login", handler.Login)

	body, _ := json.Marshal(user.LoginRequest{
		Username: "",
		Password: "",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")

	mockSvc.AssertExpectations(t)
}

func TestUserHandlers_Login_UserNotFound(t *testing.T) {
	mockSvc := mocks.NewUserService(t)
	cfg := config.Config{}
	handler := user.NewUserHandlers(mockSvc, cfg)

	r := setupRouter()
	r.POST("/api/v1/users/login", handler.Login)

	mockSvc.On("Get", mock.Anything, mock.AnythingOfType("entities.UserFilter")).Return(entities.User{}, errors.Error{})

	body, _ := json.Marshal(user.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")

	mockSvc.AssertExpectations(t)
}

func TestUserHandlers_Login_UserInactive(t *testing.T) {
	mockSvc := mocks.NewUserService(t)
	cfg := config.Config{}
	handler := user.NewUserHandlers(mockSvc, cfg)

	r := setupRouter()
	r.POST("/api/v1/users/login", handler.Login)

	userToReturn := entities.User{
		ID:     "user123",
		Status: constant.UserStatusInactive,
	}

	mockSvc.On("Get", mock.Anything, mock.AnythingOfType("entities.UserFilter")).Return(userToReturn, errors.Error{})

	body, _ := json.Marshal(user.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "user is inactive")

	mockSvc.AssertExpectations(t)
}
