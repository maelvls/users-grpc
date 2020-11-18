package grpc

import (
	"fmt"
	"net"

	service "github.com/maelvls/users-grpc/pkg/service"
	"github.com/maelvls/users-grpc/schema/user"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Run starts the server. Set reflexion to true if you want to be able to
// use grpcurl or prototool to discover the proto files.
func Run(addr string, enableReflection, tls, samples bool, certFile, keyFile string) error {
	userServer := NewUserServer()

	txn := userServer.Txn(true)
	defer userServer.Rollback(txn)

	if samples {
		logrus.Info("loading sample users, disable with --samples=false")
		err := service.LoadSampleUsers(txn)
		if err != nil {
			return fmt.Errorf("while loading sample users: %w", err)
		}
		userServer.Commit(txn)
	}

	var opts []grpc.ServerOption
	if tls {
		if certFile == "" {
			return fmt.Errorf("since --tls was given, you must also give --tls-cert-file")
		}
		if keyFile == "" {
			return fmt.Errorf("since --tls was given, you must also give --tls-key-file")
		}

		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("failed to generate TLS server credentials: %w", err)
		}

		opts = append(opts, grpc.Creds(creds))
	} else {
		logrus.Printf("TLS is disabled by default, use --tls, --tls-cert-file and --tls-key-file to enable TLS")
	}

	srv := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(srv, userServer)
	health := health.NewServer()
	health.SetServingStatus("user", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(srv, health)

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
