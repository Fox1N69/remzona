package postgres

import (
	"sso/internal/domain/models"

	"gorm.io/gorm"
)

type MigrationOption func(*migration)

func WithDebug() MigrationOption {
	return func(m *migration) {
		m.debug = true
	}
}

type Migration interface {
	AuthMigrate(...MigrationOption) error
}

type migration struct {
	db    *gorm.DB
	debug bool
}

func NewMigration(db *gorm.DB) Migration {
	return &migration{db: db}
}

func (m *migration) AuthMigrate(opts ...MigrationOption) error {
	for _, opt := range opts {
		opt(m)
	}

	if m.debug {
		m.db = m.db.Debug()
	}

	if err := m.db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}
