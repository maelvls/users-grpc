package main

import (
	"github.com/maelvls/users-grpc/pkg/cli"
)

// Set during build, e.g.: go build  -ldflags"-X main.version=$(git describe)".
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.Execute(cli.Version{
		Commit:  commit,
		Version: version,
		Date:    date,
	})
}
