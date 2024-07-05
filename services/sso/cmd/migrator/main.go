package main

import (
	"remzona-sso/infra"
	"remzona-sso/internal/models"
)

func main() {
	i := infra.New("config/config.json")
	i.SetMode()
	i.Migrate(
		&models.User{},
		&models.Login{},
		&models.Post{},
	)
}
