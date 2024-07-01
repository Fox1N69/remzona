package grpc

import (
	"fmt"
	"net"
	authgrpc "sso/internal/grpc/auth"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	log        *logrus.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	logrus *logrus.Logger,
	authService authgrpc.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        logrus,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Start gRPC server
func (a *App) Run() error {
	const pt = "grpcapp.Run"
	log := a.log.WithFields(logrus.Fields{
		"pt":   pt,
		"port": a.port,
	})

	log.Info("Starting gRPC server...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", pt, err)
	}

	log.Info("gRPC server is running", logrus.Fields{"addr": lis.Addr().String()})
	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", pt, err)
	}

	return nil
}

// Stop gRPC server
func (a *App) Stop() {
	const pt = "grpcapp.Stop"

	a.log.WithField("pt", pt).Info("stopping gRPC server", logrus.Fields{"port": a.port})

	a.gRPCServer.GracefulStop()
}
