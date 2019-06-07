package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// AppVersion is passed to cmd.Execute().
type AppVersion struct {
	Version    string // 0b5ed7a
	BuildDate  string // master
	GitCommit  string // clean or dirty
	GitBranch  string // v1.0.0, v1.0.1-5-g585c78f-dirty, fbd157c (git describe --tags --dirty --always)
	GitState   string // 2016-08-04T18:07:54Z
	GitSummary string // 2.0.0
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of this binary to stdout",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s%s\n", version.GitSummary, func() string {
			if strings.Compare(version.GitBranch, "master") == 0 {
				return ""
			}
			return fmt.Sprintf(" (%s)", version.GitBranch)
		}())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
