package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sairahul1526/morphic/config"
	"github.com/sairahul1526/morphic/logger"
	"go.uber.org/zap"
)

// request metrics

var TotalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of requests",
	},
	[]string{"path", "method", "status"},
)

var HTTPLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_response_time_ms",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{"path", "method", "status"},
)

// database metrics

var DBRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "db_requests_total",
		Help: "Number of database operation",
	},
	[]string{"model", "operation", "status", "error"},
)

var DBLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "db_response_time_ms",
		Help:    "Duration of database calls",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{"model", "operation", "status", "error"},
)

// other metrics

func Serve(ctx context.Context, cfg config.Config, db *sql.DB) {

	err := prometheus.Register(TotalRequests)
	if err != nil {
		logger.Fatal("Error registering totalRequests metric", zap.Error(err))
	}

	err = prometheus.Register(HTTPLatency)
	if err != nil {
		logger.Fatal("Error registering httpLatency metric", zap.Error(err))
	}

	err = prometheus.Register(DBRequests)
	if err != nil {
		logger.Fatal("Error registering dbRequests metric", zap.Error(err))
	}

	err = prometheus.Register(DBLatency)
	if err != nil {
		logger.Fatal("Error registering dbLatency metric", zap.Error(err))
	}

	registerSystemMetrics(cfg, db)

	// start server
	if strings.EqualFold(cfg.Metrics.Enabled, "1") {
		logger.Info("Starting metrics server")
		serveMux := &http.ServeMux{}
		server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Metrics.Port), Handler: serveMux}
		serveMux.Handle("/metrics", promhttp.Handler())

		go func() {
			err := server.ListenAndServe()
			if err != nil {
				logger.Error("Error starting metrics server", zap.Error(err))
			}
		}()

		// Listening for context cancellation to gracefully shut down the server
		go func() {
			<-ctx.Done()
			if err := server.Close(); err != nil {
				logger.Fatal("Metrics server Close:", zap.Error(err))
			}
		}()
	}
}

func registerSystemMetrics(cfg config.Config, db *sql.DB) {
	// Note: not adding GoCollector and ProcessCollector as they are included by default in the Default registry
	tokens := strings.Split(cfg.Database.URL, "/")
	rawDbName := tokens[len(tokens)-1]
	dbName := strings.Split(rawDbName, "?")[0]
	err := prometheus.Register(collectors.NewDBStatsCollector(db, dbName))
	if err != nil {
		logger.Fatal("Error registering dbStats collector", zap.Error(err))
	}
}
