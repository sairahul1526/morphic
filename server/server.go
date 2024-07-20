package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/sairahul1526/morphic/metrics"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sairahul1526/morphic/config"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/store"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/sairahul1526/morphic/docs"
)

// Serve initialises all the HTTP API routes, starts listening for requests at addr, and blocks until
// server exits. Server exits gracefully when context is cancelled.
func Serve(ctx context.Context, addr string, cfg config.Config) error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(
		RequestID(),
		gin.Recovery(),
		gzip.Gzip(gzip.DefaultCompression),
		PrometheusMiddleware(),
		DefaultStructuredLogger(),
		Authorize(cfg.Auth.Secret),
	)

	// address for swagger docs
	docs.SwaggerInfo.Host = strings.ReplaceAll(addr, "0.0.0.0", "localhost")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(SetJSONContentType())

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Service.AllowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowHeaders: []string{"*"},
	}))

	db, err := store.Client(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", []zap.Field{
			zap.Error(err),
		}...)
		return err
	}

	// start metrics server
	metrics.Serve(ctx, cfg, db.DB.DB)

	registerEmployeeV1APIs(db, router.Group("/api/v1"))
	registerUserV1APIs(db, router.Group("/api/v1"), cfg)

	logger.Info("starting server", zap.String("addr", addr))

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Listening for context cancellation to gracefully shut down the server
	go func() {
		<-ctx.Done()
		if err := server.Close(); err != nil {
			logger.Fatal("API server Close:", zap.Error(err))
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("API server closed under request")
		} else {
			logger.Fatal("API server closed unexpectedly:", zap.Error(err))
		}
	}
	return err
}
