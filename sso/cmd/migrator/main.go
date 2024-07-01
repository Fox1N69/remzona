package main

import (
	"sso/internal/config"
	"sso/storage/postgres"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.MustLoad()

	db, err := postgres.New("user=postgres password=8008 dbname=grpc_remzona port=5432 sslmode=disable")
	if err != nil {
		logrus.Fatal("Error connect to database: ", err)
		return
	}

	if !cfg.Debug {
		if err := db.Migration.AuthMigrate(); err != nil {
			logrus.Error("Error migrate table for auth service", err)
			return
		}
	} else {
		if err := db.Migration.AuthMigrate(postgres.WithDebug()); err != nil {
			logrus.Error("Error migrate table for auth service", err)
			return
		}
	}

	logrus.Info("Database migrate is successfully")
}
