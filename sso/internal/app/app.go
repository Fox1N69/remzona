package app

import (
	"time"

	grpcapp "sso/internal/app/grpc"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *logrus.Logger, grpcPort int, storagePath string, TokenTTL time.Duration) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
