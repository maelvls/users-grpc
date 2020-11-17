package grpc

import (
	"fmt"
	"net"

	service "github.com/maelvls/users-grpc/pkg/service"
	"github.com/maelvls/users-grpc/schema/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Run starts the server.
func Run(addr string) error {
	server := NewUserServer()

	txn := server.Txn(true)
	defer server.Rollback(txn)

	err := service.LoadSampleUsers(txn)
	if err != nil {
		return fmt.Errorf("while loading sample users: %w", err)
	}
	server.Commit(txn)

	srv := grpc.NewServer()
	user.RegisterUserServiceServer(srv, server)
	grpc_health_v1.RegisterHealthServer(srv, &health.Server{})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	return srv.Serve(lis)
}
