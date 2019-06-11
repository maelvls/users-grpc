package service

import (
	"context"

	v1 "github.com/maelvls/users-grpc/schema/health/v1"
)

// HealthImpl implements the HealthServer interface (see health.pb.go file).
type HealthImpl struct{}

// Check returns SERVING.
func (h *HealthImpl) Check(ctx context.Context, args *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{
		Status: v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is not implemented for now.
func (h *HealthImpl) Watch(*v1.HealthCheckRequest, v1.Health_WatchServer) error {
	return nil
}
