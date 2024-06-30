package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email"`
	PassHash []byte `json:"pass_hash"`
}
