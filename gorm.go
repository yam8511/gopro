package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User 使用者型態
type User struct {
	gorm.Model
	Name string
}

func newDBConnection() (db *gorm.DB, err error) {
	connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
	)
	fmt.Println("Connection ---> ", connection)
	db, err = gorm.Open("postgres", connection)
	if err != nil {
		return
	}
	err = db.LogMode(true).Error
	if err != nil {
		return
	}
	err = db.AutoMigrate(new(User)).Error
	return
}
