package storage

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitGormDB() *gorm.DB {
	dsn := "user=postgres password=8008 dbname=grpc_remzona port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Error connect to database: ", err)
		return nil
	}
	logrus.Info("Server connet to database: " + "grpc_remzona")

	DB = db

	return DB
}