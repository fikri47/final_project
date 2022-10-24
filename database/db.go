package database

import (
	"final_project/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "12345"
	dbPort   = 5432
	dbName   = "db_mygram"
	db       *gorm.DB
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		fmt.Println("error open connection to database", err.Error())
		return
	}

	err = db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})

	if err != nil {
		fmt.Println("error migration", err.Error())
		return
	}

	fmt.Println("succesfully connected to database")
}

func GetDB() *gorm.DB {
	return db
}
