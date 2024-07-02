package app

import (
	"time"

	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/storage/postgres"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *logrus.Logger, grpcPort int, TokenTTL time.Duration) *App {
	storage, err := postgres.New("user=postgres password=8008 dbname=grpc_remzona port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, TokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
