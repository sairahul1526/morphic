package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/raystack/salt/db"
	"github.com/rs/xid"
	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/metrics"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const (
	headerRequestID = "X-Request-Id"
)

func requestLogFields(c *gin.Context) []zap.Field {
	clientID, _, _ := c.Request.BasicAuth()

	return []zap.Field{
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("request_id", c.Request.Header.Get(headerRequestID)),
		zap.String("client_id", clientID),
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := strings.TrimSpace(c.Request.Header.Get(headerRequestID))
		if rid == "" {
			rid = xid.New().String()
		}

		headers := c.Request.Header.Clone()
		headers.Set(headerRequestID, rid)

		c.Writer.Header().Set(headerRequestID, rid)
		c.Request.Header = headers
		c.Next()
	}
}

func Authorize(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				logger.Error("No credentials supplied", append(requestLogFields(c), []zap.Field{zap.String("authHeader", authHeader)}...)...)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No credentials supplied"})
				return
			}

			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				// Validate the alg is what you expect
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Return the secret key for validation
				return []byte(secret), nil
			})

			if err != nil || !parsedToken.Valid {
				logger.Error("Unauthorized", append(requestLogFields(c), []zap.Field{zap.Error(err)}...)...)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				logger.Error("Unauthorized", append(requestLogFields(c), []zap.Field{zap.Error(err)}...)...)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			// Check token expiration
			exp, ok := claims["exp"].(float64)
			if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
				logger.Error("Token expired", append(requestLogFields(c), []zap.Field{zap.String("error", "token expired")}...)...)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				logger.Error("Unauthorized", append(requestLogFields(c), []zap.Field{zap.String("error", "user_id not found in token claims")}...)...)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			headers := c.Request.Header.Clone()
			headers.Set("user_id", userID)
			c.Request.Header = headers
			c.Next()
		} else if c.FullPath() == "/api/v1/users/login" || strings.Contains(c.FullPath(), "swagger") || strings.Contains(c.FullPath(), "ping") {
			c.Next()
		} else {
			logger.Error("No credentials supplied", append(requestLogFields(c), []zap.Field{zap.String("authHeader", authHeader)}...)...)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No credentials supplied"})
			return
		}
	}
}

func SetJSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		method := c.Request.Method
		path := c.FullPath()

		c.Next() // Process the request
		statusCode := c.Writer.Status()

		// Use the generalizePath function to get the handler name
		handlerName := generalizePath(method, path)

		metrics.TotalRequests.WithLabelValues(handlerName, method, strconv.Itoa(statusCode)).Inc()
		metrics.HTTPLatency.WithLabelValues(handlerName, method, strconv.Itoa(statusCode)).Observe(float64(time.Since(startTime).Milliseconds()))
	}
}

func generalizePath(method, path string) string {
	// use this to generalize the path, remove ids from the path and replace them with a readble name
	switch {
	case strings.Contains(path, "/api/v1/employees"):
		return "/api/v1/employees"
	case strings.Contains(path, "/api/v1/users"):
		return "/api/v1/users"
	default:
		// Handle other routes or return a default value
		return "Unknown"
	}
}

// DefaultStructuredLogger logs a gin HTTP request in JSON format. Uses the
// default logger from rs/zerolog.
func DefaultStructuredLogger() gin.HandlerFunc {
	return StructuredLogger(logger.GetLogger())
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// StructuredLogger logs a gin HTTP request in JSON format
func StructuredLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var bodyBytes []byte
		// Read the body
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Wrap the response writer to capture the response body
		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		if c.Writer.Status() >= 400 {
			logger.Error(param.ErrorMessage, []zap.Field{
				zap.String("client_id", param.ClientIP),
				zap.String("method", param.Method),
				zap.Int("status_code", param.StatusCode),
				zap.Int("body_size", param.BodySize),
				zap.String("path", param.Path),
				zap.String("latency", param.Latency.String()),
				zap.String("request body", string(bodyBytes)),
				zap.String("response body", w.body.String()),
				zap.Error(errors.New(param.ErrorMessage)),
			}...)
		} else {
			logger.Info("", []zap.Field{
				zap.String("client_id", param.ClientIP),
				zap.String("method", param.Method),
				zap.Int("status_code", param.StatusCode),
				zap.Int("body_size", param.BodySize),
				zap.String("path", param.Path),
				zap.String("latency", param.Latency.String()),
			}...)
		}
	}
}

// DBTransactionMiddleware : to setup the database transaction middleware
func DBTransactionMiddleware(db *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		txHandle, err := db.BeginTxx(c, nil)
		if err != nil {
			logger.Error("Failed to begin transaction", []zap.Field{
				zap.Error(err),
			}...)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		logger.Debug("beginning database transaction", []zap.Field{
			zap.String("request_id", c.Request.Header.Get(headerRequestID)),
		}...)

		defer func() {
			if r := recover(); r != nil {
				if err = txHandle.Rollback(); err != nil && !strings.Contains(err.Error(), "transaction has already been committed or rolled back") {
					logger.Error("failed to rollback transaction", []zap.Field{
						zap.Error(err),
					}...)
				}
			}
		}()

		c.Set(constant.DBTrxKey, txHandle)
		c.Next()

		if c.Writer.Status() < 300 {
			logger.Debug("committing transactions", []zap.Field{
				zap.String("request_id", c.Request.Header.Get(headerRequestID)),
			}...)
			if err := txHandle.Commit(); err != nil && !strings.Contains(err.Error(), "transaction has already been committed or rolled back") {
				logger.Error("Failed to commit transaction", []zap.Field{
					zap.Error(err),
				}...)
			}
		} else {
			logger.Debug("rolling back transaction", []zap.Field{
				zap.String("request_id", c.Request.Header.Get(headerRequestID)),
			}...)
			err = txHandle.Rollback()
			if err != nil && !strings.Contains(err.Error(), "transaction has already been committed or rolled back") {
				logger.Error("failed to rollback transaction", []zap.Field{
					zap.Error(err),
				}...)
			}
		}
	}
}
