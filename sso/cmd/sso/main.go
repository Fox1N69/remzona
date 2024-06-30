package main

import (
	"sso/internal/app"
	"sso/internal/config"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.MustLoad()

	log := logrus.New()

	log.Info("starting app", logrus.Fields{"config": cfg})

	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTl)

	application.GRPCSrv.MustRun()
}
