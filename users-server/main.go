package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	grpc_health_v1 "github.com/maelvls/users-grpc/schema/health/v1"
	"github.com/maelvls/users-grpc/schema/user"
	"github.com/maelvls/users-grpc/users-server/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Set during build, e.g.: go build  -ldflags"-X main.version=$(git
// describe)". Some vars are commented out so that 'golangci-lint run' us
// happy, but you can un-comment them and use them.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	addr    = flag.String("address", ":8000", "Address used by the server to start listening. Default is ':8000' (equivalent to 0.0.0.0:8000).")
	logfmt  = flag.String("logfmt", "text", "Log format ('text', 'json'). Default is 'text'.")
	verbose = flag.Bool("v", false, "Make the server more verbose.")
)

func main() {
	flag.Parse()
	// Set the log format according to the --logfmt flag.
	switch *logfmt {
	case "", "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Errorf(`valid values for --logfmt are 'json' or 'text'`)
		os.Exit(1)
	}

	// Set to verbose.
	if *verbose {
		log.SetLevel(log.TraceLevel)
	}

	log.Printf("listening on address %s, version %s (git %s, built on %s)", *addr, version, commit, date)

	if err := Run(*addr); err != nil {
		log.Errorf("launching server: %v", err)
		os.Exit(1)
	}
}

// Run starts the server.
func Run(addr string) error {
	svc := service.NewUserImpl()
	err := svc.LoadSampleUsers()
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	user.RegisterUserServiceServer(srv, svc)
	grpc_health_v1.RegisterHealthServer(srv, &service.HealthImpl{})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	return srv.Serve(lis)
}