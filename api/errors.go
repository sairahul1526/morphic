package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sairahul1526/morphic/pkg/errors"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}
type Error struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func HarmonizeErrorResponse(c *gin.Context, err errors.Error) {
	switch err.Cause {
	case errors.ErrCodeBadRequest:
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: Error{Code: err.Code, Message: err.Message, Details: err.Details}})
		return
	case errors.ErrCodeNotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: Error{Code: err.Code, Message: err.Message, Details: err.Details}})
		return
	case errors.ErrCodeInternalServer:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: Error{Code: err.Code, Message: err.Message, Details: err.Details}})
		return
	case errors.ErrForbidden:
		c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: Error{Code: err.Code, Message: err.Message, Details: err.Details}})
		return
	case errors.ErrConflictCode:
		c.AbortWithStatusJSON(http.StatusConflict, ErrorResponse{Error: Error{Code: err.Code, Message: err.Message, Details: err.Details}})
		return
	case errors.ErrUnsupportedCode:
		c.AbortWithStatusJSON(http.StatusNotImplemented, ErrorResponse{Error: Error{Code: err.Code, Message: err.Message, Details: err.Details}})
		return
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: Error{Code: err.Code, Message: err.Message, Details: err.Details}})
	}
}
