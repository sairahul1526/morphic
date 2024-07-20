package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/raystack/salt/db"
	"github.com/sairahul1526/morphic/constant"
	"github.com/stretchr/testify/assert"
)

func MockClient(mockDB *sql.DB) *db.Client {
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return &db.Client{DB: sqlxDB}
}

func TestDBTransactionMiddleware_Success(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection")
	defer db.Close()

	sqlxDB := MockClient(db)

	dbmock.ExpectBegin()
	dbmock.ExpectCommit()

	r := gin.Default()
	r.Use(DBTransactionMiddleware(sqlxDB))
	r.GET("/test", func(c *gin.Context) {
		txHandle, exists := c.Get(constant.DBTrxKey)
		assert.True(t, exists)
		assert.NotNil(t, txHandle)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, dbmock.ExpectationsWereMet(), "there were unfulfilled expectations")

	assert.Nil(t, dbmock.ExpectationsWereMet())
}

func TestDBTransactionMiddleware_FailureToBeginTransaction(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection")
	defer db.Close()

	sqlxDB := MockClient(db)

	dbmock.ExpectBegin().WillReturnError(fmt.Errorf("begin transaction error"))

	r := gin.Default()
	r.Use(DBTransactionMiddleware(sqlxDB))
	r.GET("/test", func(c *gin.Context) {
		// This should not be reached due to begin transaction failure
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, dbmock.ExpectationsWereMet(), "there were unfulfilled expectations")

	assert.Nil(t, dbmock.ExpectationsWereMet())
}

func TestDBTransactionMiddleware_RollbackOnFailure(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection")
	defer db.Close()

	sqlxDB := MockClient(db)

	dbmock.ExpectBegin()
	dbmock.ExpectRollback()

	r := gin.Default()
	r.Use(DBTransactionMiddleware(sqlxDB))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NoError(t, dbmock.ExpectationsWereMet(), "there were unfulfilled expectations")

	assert.Nil(t, dbmock.ExpectationsWereMet())
}
