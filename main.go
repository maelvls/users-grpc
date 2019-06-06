package main

import (
	"github.com/maelvalais/quote/cmd"
)

// These fields are populated by github.com/ahmetb/govvv. These variables
// must stay in the main package.
var (
	Version    = "untouched"
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
)

func main() {
	cmd.Execute(cmd.AppVersion{
		Version:    Version,
		BuildDate:  BuildDate,
		GitCommit:  GitCommit,
		GitBranch:  GitBranch,
		GitState:   GitState,
		GitSummary: GitSummary,
	})
}
