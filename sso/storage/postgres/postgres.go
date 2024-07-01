package postgres

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/storage"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
	Migration
}

// New - fucntion for initialization postgres database with gorm.io
//
// return gorm.DB and Migration or errors
func New(dsn string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Error connect to database: ", err)
		return nil, err
	}
	logrus.Info("Server connet to database: " + "grpc_remzona")

	return &Storage{
		db:        db,
		Migration: NewMigration(db),
	}, nil
}

// Stop - furnction for stop connect to database
func (s *Storage) Stop() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uint64, error) {
	const op = "storage.postgres.SaveUser"

	user := models.User{
		Email:    email,
		PassHash: passHash,
	}

	err := s.db.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		} else if errors.Is(err, gorm.ErrInvalidData) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidData)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uint64(user.ID), nil
}
