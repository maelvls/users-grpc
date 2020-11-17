package grpc

import (
	"fmt"
	"net"

	service "github.com/maelvls/users-grpc/pkg/service"
	"github.com/maelvls/users-grpc/schema/user"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Run starts the server. Set reflexion to true if you want to be able to
// use grpcurl or prototool to discover the proto files.
func Run(addr string, enableReflection bool) error {
	userServer := NewUserServer()

	txn := userServer.Txn(true)
	defer userServer.Rollback(txn)

	err := service.LoadSampleUsers(txn)
	if err != nil {
		return fmt.Errorf("while loading sample users: %w", err)
	}
	userServer.Commit(txn)

	srv := grpc.NewServer()
	user.RegisterUserServiceServer(srv, userServer)
	grpc_health_v1.RegisterHealthServer(srv, &health.Server{})

	if enableReflection {
		logrus.Info("reflection enabled, you can now use tools like grpcurl")
		reflection.Register(srv)
	} else {
		logrus.Info("reflection disabled by default, use --reflection to be able to use tools like grpcurl")
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	return srv.Serve(lis)
}
