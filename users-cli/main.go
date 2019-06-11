package main

import (
	"github.com/maelvls/users-grpc/users-cli/cmd"
)

// Set during build, e.g.: go build  -ldflags"-X main.version=$(git describe)".
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(cmd.Version{
		Commit:  commit,
		Version: version,
		Date:    date,
	})
}
