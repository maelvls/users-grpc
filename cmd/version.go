package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// AppVersion is passed to cmd.Execute().
type AppVersion struct {
	Version    string
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
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
