package database

import (
	"UnitUser/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
/// testing
/// testing two

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
func CreateDB() {
	DSN := "host=localhost user=postgres password=123 dbname=usertest port=5432"
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	AutoMigrate(db)
	DB = db
}
func SetDB(database *gorm.DB) {
	DB = database
}
