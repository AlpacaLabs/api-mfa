package grpc

import (
	"context"

	"github.com/sirupsen/logrus"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

func (s Server) Check(ctx context.Context, request *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// rpc Watch(HealthCheckRequest) returns (stream HealthCheckResponse);
func (s Server) Watch(request *health.HealthCheckRequest, stream health.Health_WatchServer) error {
	if err := stream.Send(&health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}); err != nil {
		logrus.Errorf("server failed to stream HealthCheckResponse: %v", err)
	}

	return nil
}
