package models

import (
	"database/sql"
	"time"
)

type GormCustom struct {
	ID        uint64       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" sql:"index"`
}

type User struct {
	GormCustom
	Email    string `json:"email"`
	PassHash []byte `json:"pass_hash"`
	IsAdmin  bool   `json:"is_admin"`
}
