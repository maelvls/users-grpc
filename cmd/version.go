package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AppVersion is passed to cmd.Execute().
type AppVersion struct {
	BuildDate  string // 2016-08-04T18:07:54Z
	GitSummary string // v1.0.0, v1.0.1-5-g585c78f-dirty, fbd157c (git describe --tags --dirty --always)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of this binary to stdout",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", version.GitSummary)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
