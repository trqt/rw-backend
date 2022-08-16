package database

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"readyworker.com/backend/model"
)

func Connect() (db *gorm.DB) {
	dbUser := "readyworker"
	dbPass := os.Getenv("MARIADB_PASSWORD")
	dbName := "ReadyWorker"
	dsn := dbUser + ":" + dbPass + "@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&model.User{})

	return db
}
