package main

import (
	"github.com/maelvls/quote/cmd"
)

// These fields are populated by github.com/ahmetb/govvv. These variables
// must stay in the main package.
var (
	BuildDate  string
	GitSummary string
)

func main() {
	cmd.Execute(cmd.AppVersion{
		BuildDate:  BuildDate,
		GitSummary: GitSummary,
	})
}
