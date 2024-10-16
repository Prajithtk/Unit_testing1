package test

import (
	"UnitUser/controller"
	"UnitUser/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDBSignup() (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic("falied to open mock database connection")
	}
	dialector := postgres.New(postgres.Config{
		Conn: mockDB,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to open gorm db connection")
	}
	// new change implemented
	database.SetDB(gormDB)
	return gormDB, mock
}
func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("a valid signup", func(t *testing.T) {
		testDB, mock := SetupTestDBSignup()
		sqlDB, err := testDB.DB()
		require.NoError(t, err)
		defer sqlDB.Close()

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users" \("created_at","updated_at","deleted_at","name","email","password"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Test User", "testuser@example.com", sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		router := gin.Default()
		router.POST("/signup", controller.UserSignUp)

		payload := `{"name": "Test User", "email": "testuser@example.com", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		w:= httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK{
			t.Logf("Response Code: %d", w.Code)
			t.Logf("Response Body: %s", w.Body.String())
		}

		require.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "user created successfully.",response["message"])

		if err := mock.ExpectationsWereMet(); err != nil{
			t.Errorf("mock expectations were not met: %v", err)
		}
	})

	t.Run("invalid signup", func(t *testing.T) {
		router := gin.Default()
		router.POST("/signup", controller.UserSignUp)

		payload := `{"name": "", "email": "", "password": ""}`
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Logf("Response Code: %d", w.Code)
			t.Logf("Response Body: %s", w.Body.String())
		}

		require.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "all fields are mandatory", response["error"])
	})
}
