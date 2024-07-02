package models

import "gorm.io/gorm"

type App struct {
	gorm.Model
	Name   string
	Secret string
}
