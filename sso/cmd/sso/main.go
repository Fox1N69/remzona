package main

import (
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.MustLoad()

	log := logrus.New()

	log.Info("starting app", logrus.Fields{"config": cfg})

	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTl)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()

	log.Infoln("Server stopped")
}
