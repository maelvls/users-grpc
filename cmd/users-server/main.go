package main

import (
	"flag"
	"os"

	grpc "github.com/maelvls/users-grpc/pkg/grpc"
	"github.com/sirupsen/logrus"
)

// Set during build, e.g.: go build  -ldflags"-X main.version=$(git
// describe)". Some vars are commented out so that 'golangci-lint run' us
// happy, but you can un-comment them and use them.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	addr     = flag.String("address", ":8000", "Address used by the server to start listening. Default is ':8000' (equivalent to 0.0.0.0:8000).")
	logfmt   = flag.String("logfmt", "text", "Log format ('text', 'json'). Default is 'text'.")
	verbose  = flag.Bool("v", false, "Make the server more verbose.")
	tls      = flag.Bool("tls", false, "If set, the connection is established with TLS; otherwise, connection is in clear text mode (h2c, HTTP/2 clear text).")
	certFile = flag.String("tls-cert-file", "", "The TLS cert file, required if --tls is set.")
	keyFile  = flag.String("tls-key-file", "", "The TLS key file, required if --tls is set.")
	samples  = flag.Bool("samples", false, "Load some user samples on startup.")
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md
	reflection = flag.Bool("reflection", false, "Enable reflection, useful for using grpcurl or related tools.")
)

func main() {
	flag.Parse()
	// Set the log format according to the --logfmt flag.
	switch *logfmt {
	case "", "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.Errorf(`valid values for --logfmt are 'json' or 'text'`)
		os.Exit(1)
	}

	// Set to verbose.
	if *verbose {
		logrus.SetLevel(logrus.TraceLevel)
	}

	logrus.Printf("listening on address %s, version %s (git %s, built on %s)", *addr, version, commit, date)

	if err := grpc.Run(*addr, *reflection, *tls, *samples, *certFile, *keyFile); err != nil {
		logrus.Errorf("launching server: %v", err)
		os.Exit(1)
	}
}
