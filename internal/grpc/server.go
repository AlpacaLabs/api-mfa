package grpc

import (
	"fmt"
	"net"

	"github.com/AlpacaLabs/api-mfa/internal/configuration"
	"github.com/AlpacaLabs/api-mfa/internal/service"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config  configuration.Config
	service service.Service
}

func NewServer(config configuration.Config, service service.Service) Server {
	return Server{
		config:  config,
		service: service,
	}
}

func (s Server) Run() {
	address := fmt.Sprintf(":%d", s.config.GrpcPort)

	log.Infof("Preparing to listen for gRPC on %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register our services
	mfaV1.RegisterMFAServiceServer(grpcServer, s)
	health.RegisterHealthServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Debug("Registered gRPC services...")

	log.Infof("Starting gRPC server on %s...", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
