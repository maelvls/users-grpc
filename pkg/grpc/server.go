package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	service "github.com/maelvls/users-grpc/pkg/service"
	"github.com/maelvls/users-grpc/schema/user"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Run starts the server. Set reflexion to true if you want to be able to
// use grpcurl or prototool to discover the proto files.
func Run(ctx context.Context, addr, addrMetrics string, enableReflection, tls, samples bool, certFile, keyFile string) error {
	userServer := NewUserServer()

	if samples {
		logrus.Info("loading sample users, disable with --samples=false")

		txn := userServer.Txn(true)
		defer userServer.Rollback(txn)
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

	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_prometheus.UnaryServerInterceptor,
		grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel)),
	)))

	srv := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(srv, userServer)
	health := health.NewServer()
	health.SetServingStatus("user", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(srv, health)
	grpc_prometheus.Register(srv)

	if enableReflection {
		logrus.Info("reflection enabled, you can now use tools like grpcurl")
		reflection.Register(srv)
	} else {
		logrus.Info("reflection disabled by default, use --reflection to be able to use tools like grpcurl")
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	setupSignalHandler(cancel)

	group, _ := errgroup.WithContext(ctx)

	metrics := &http.Server{Addr: addrMetrics, Handler: promhttp.Handler()}
	group.Go(func() error {
		defer cancel()
		return metrics.ListenAndServe()
	})

	group.Go(func() error {
		defer cancel()
		return srv.Serve(lis)
	})

	group.Go(func() error {
		// Cleanup goroutine.
		<-ctx.Done()
		srv.GracefulStop()
		_ = metrics.Shutdown(context.Background())
		return nil
	})

	return group.Wait()
}

// setupSignalHandler will call handleShutdown as soon as SIGINT or SIGTERM
// is caught. If a second signal is received afterwards, the program exits
// immediatly.
//
// It must only be called once.
func setupSignalHandler(handleShutdown func()) {
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		s := <-c
		logrus.Infof("Signal %s received, shutting down gracefully...", s.String())
		handleShutdown()
		s = <-c
		logrus.Infof("Signal %s received again: aborting graceful shutdown", s.String())
		os.Exit(1)
	}()
}
