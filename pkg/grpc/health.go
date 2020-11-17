package grpc

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthImpl implements the HealthServer interface (see health.pb.go file).
type HealthImpl struct{}

// Check returns SERVING.
func (h *HealthImpl) Check(ctx context.Context, args *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is not implemented for now.
func (h *HealthImpl) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}
