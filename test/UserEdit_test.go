package test

import (
	"UnitUser/controller"
	"UnitUser/database"
	"UnitUser/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func UserEditMock() (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to open sqlmock database connection")
	}
	dialecter := postgres.New(postgres.Config{
		Conn: mockDB,
	})
	gormDB, err := gorm.Open(dialecter, &gorm.Config{})
	if err != nil {
		panic("failed to open gorm db connection")
	}
	database.SetDB(gormDB)

	users := []models.User{
		{Model: gorm.Model{ID: 1}, Name: "user1", Email: "user1@gmail.com", Password: "user1@123"},
		{Model: gorm.Model{ID: 2}, Name: "user2", Email: "user2@gmail.com", Password: "user2@123"},
	}
	mock.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(0, 0))
	database.DB.Exec("DELETE FROM users")
	if err := database.DB.Create(&users).Error; err != nil {
		fmt.Println("-----------", err)
	}

	cleanup := func() {
		mockDB.Close()
	}
	return mock, cleanup
}

func TestUserEdit(t *testing.T) {
	mock, cleanup := UserEditMock()
	defer cleanup()

	t.Run("successful edit user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "email"=\$1,"name"=\$2,"password"=\$3,"updated_at"=\$4 WHERE id = \$5 AND "users"."deleted_at" IS NULL`).
			WithArgs("userEdit@gmail.com", "userEdit", "userEdit@123", sqlmock.AnyArg(), 11).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		router := gin.Default()
		router.PATCH("/user/edit/:id", controller.EditUser)
/// tiy
		user := models.User{
			Name:     "userEdit",
			Email:    "userEdit@gmail.com",
			Password: "userEdit@123",
		}
		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPatch, "/user/edit/11", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		fmt.Println("=====================",w.Code)
		assert.Contains(t, w.Body.String(), "Successfully updated user")
	})
}
