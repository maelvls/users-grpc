package grpc

import (
	"fmt"
	"net"

	service "github.com/maelvls/users-grpc/pkg/service"
	"github.com/maelvls/users-grpc/schema/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Run starts the server.
func Run(addr string) error {
	svc := NewUserImpl()

	txn := svc.DB.Txn(true)
	defer txn.Abort()

	err := service.LoadSampleUsers(txn)
	if err != nil {
		return fmt.Errorf("while loading sample users: %w", err)
	}
	txn.Commit()

	srv := grpc.NewServer()
	user.RegisterUserServiceServer(srv, svc)
	grpc_health_v1.RegisterHealthServer(srv, &HealthImpl{})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	return srv.Serve(lis)
}
