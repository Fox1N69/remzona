package main

import (
	"remzona-sso/infra"
	"remzona-sso/internal/api"
)

func main() {
	i := infra.New("config/config.json")
	i.SetMode()

	redisClient := i.RedisClient()

	api.NewServer(i, redisClient).Run()
}
