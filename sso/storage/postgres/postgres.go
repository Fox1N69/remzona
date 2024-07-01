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

// SaveUser - save user to database
// and return userID or error
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uint64, error) {
	const op = "storage.postgres.SaveUser"

	user := models.User{
		Email:    email,
		PassHash: passHash,
	}

	err := s.db.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("%w: %s", op, storage.ErrUserNotFound)
		} else if errors.Is(err, gorm.ErrInvalidData) {
			return 0, fmt.Errorf("%w: %s", op, storage.ErrInvalidData)
		}

		return 0, fmt.Errorf("%w: %s", op, err)
	}

	return user.ID, nil
}

// User return user by email or error
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User

	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("%w: %s", op, storage.ErrUserNotFound)
		} else if errors.Is(err, gorm.ErrInvalidData) {
			return models.User{}, fmt.Errorf("%w: %s", op, storage.ErrInvalidData)
		}

		return models.User{}, fmt.Errorf("%w: %s", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var user models.User

	if err := s.db.Where("id = ?", userID).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("%w: %s", op, storage.ErrUserNotFound)
		} else if errors.Is(err, gorm.ErrInvalidData) {
			return false, fmt.Errorf("%w: %s", op, storage.ErrInvalidData)
		}

		return false, fmt.Errorf("%w: %s", op, err)
	}

	return user.IsAdmin, nil
}

// App return app or error
func (s *Storage) App(ctx context.Context, appID uint64) (models.App, error) {
	const op = "storage.postgres.App"

	var app models.App

	if err := s.db.Where("id = ?", appID).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.App{}, fmt.Errorf("%w: %s", op, storage.ErrUserNotFound)
		} else if errors.Is(err, gorm.ErrInvalidData) {
			return models.App{}, fmt.Errorf("%w: %s", op, storage.ErrInvalidData)
		}

		return models.App{}, fmt.Errorf("%w: %s", op, err)
	}

	return app, nil
}
