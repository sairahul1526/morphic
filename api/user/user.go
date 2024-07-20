package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/sairahul1526/morphic/api"
	"github.com/sairahul1526/morphic/config"
	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/pkg/auth"
	"github.com/sairahul1526/morphic/pkg/errors"
	"github.com/sairahul1526/morphic/pkg/validator"
	"go.uber.org/zap"
)

type UserService interface {
	Get(ctx context.Context, filters entities.UserFilter) (entities.User, errors.Error)
}

type UserHandlers struct {
	svc UserService
	cfg config.Config
}

func NewUserHandlers(svc UserService, cfg config.Config) *UserHandlers {
	return &UserHandlers{svc, cfg}
}

// @Tags User
// @Summary Login a user
// @Router /api/v1/users/login [post]
// @Param body body LoginRequest true "Request Body"
// @Success 200 {object} ReadResponse
// @Accept json
// @Produce json
// @failure 400 {object} api.ErrorResponse
// @failure 404 {object} api.ErrorResponse
// @failure 403 {object} api.ErrorResponse
// @failure 409 {object} api.ErrorResponse
// @failure 500 {object} api.ErrorResponse
// @failure 501 {object} api.ErrorResponse
func (h *UserHandlers) Login(c *gin.Context) {
	var body LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("failed to bind request body", zap.Error(err))
		api.HarmonizeErrorResponse(c, errors.Error{Cause: errors.ErrCodeBadRequest, Message: err.Error(), Code: errors.ErrCodeBadRequest})
		return
	}

	validate := validator.GetValidator()
	if validationErr := validate.Struct(body); validationErr != nil {
		logger.Error("failed to login", []zap.Field{
			zap.Any("user", body),
			zap.Error(validationErr),
		}...)
		api.HarmonizeErrorResponse(c, errors.Error{Cause: errors.ErrCodeBadRequest, Message: validationErr.Error(), Code: errors.ErrCodeInvalidRequest})
		return
	}

	body.Password = auth.HashPassword(body.Password)

	user, err := h.svc.Get(c, body.toDomain())
	if !err.IsEmpty() {
		logger.Error("failed to login", []zap.Field{
			zap.Any("user", body),
			zap.Error(err),
		}...)
		api.HarmonizeErrorResponse(c, err)
		return
	}

	if len(user.ID) == 0 {
		api.HarmonizeErrorResponse(c, errors.Error{Cause: errors.ErrCodeNotFound, Message: "user not found", Code: errors.ErrCodeNotFound})
		return
	}
	if user.Status == constant.UserStatusInactive {
		api.HarmonizeErrorResponse(c, errors.Error{Cause: errors.ErrForbidden, Message: "user is inactive", Code: errors.ErrForbidden})
		return
	}

	token, authErr := auth.CreateAccessToken(user.ID, h.cfg.Auth.Secret, h.cfg.Auth.Expiry)
	if authErr != nil {
		logger.Error("failed to generate token", zap.Error(authErr))
		api.HarmonizeErrorResponse(c, errors.Error{Cause: errors.ErrCodeInternalServer, Message: "failed to generate token", Code: errors.ErrCodeInternalServer})
		return
	}

	c.Header("access_token", token)

	c.JSON(http.StatusOK, NewReadResponse(user))
}
