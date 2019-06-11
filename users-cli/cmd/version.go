package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is passed to cmd.Execute().
type Version struct {
	Date    string
	Version string
	Commit  string
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and git commit to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s)\n", version.Version, version.Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
