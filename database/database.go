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
	dsn := dbUser + ":" + dbPass + "@tcp(mariadb:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Comment{})
}

func Disconnect(db *gorm.DB) error {
	sqldb, err := db.DB()
	if err != nil {
		return err
	}
	err = sqldb.Close()
	return err
}
