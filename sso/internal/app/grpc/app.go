package grpc

import (
	authgrpc "sso/internal/grpc/auth"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	log        *logrus.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logrus *logrus.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log:        logrus,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
