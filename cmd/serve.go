package cmd

import (
	"fmt"
	"log"

	quote "github.com/maelvls/quote/services/quote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		port, addr := viper.GetInt("port"), viper.GetString("addr")
		fmt.Printf("serving on %v:%v\n", port, addr)
		s := quote.NewServer()
		if err := s.Run(addr, port); err != nil {
			log.Fatalf("launching server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 8040, "Port to run quote server on; alternatively, use PORT var")
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	serveCmd.Flags().String("addr", "", "Address the server will bind to; alternatively, use ADDR var")
	viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))

	// serveCmd.Flags().String("logLevel", "DEBUG", "Log level: DEBUG, INFO, WARN, ERROR")
	// viper.BindPFlag("logLevel", serveCmd.Flags().Lookup("logLevel"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
