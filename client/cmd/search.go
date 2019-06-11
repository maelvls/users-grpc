package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/maelvls/quote/schema/user"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "searchs users from the remote quote service",
		Run: func(searchCmd *cobra.Command, args []string) {
			addr := viper.GetString("address")
			cc, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "grpc client: %v\n", err)
				os.Exit(1)
			}

			client := user.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			resp, err := client.SearchName(ctx, &user.SearchNameReq{Query: ""})

			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "grpc client: %v\n", err)
				os.Exit(1)
			}

			if resp.GetStatus().GetCode() != user.Status_SUCCESS {
				_, _ = fmt.Fprintf(os.Stderr, "grpc client: %v\n", resp.GetStatus())
				os.Exit(1)
			}

			for v := range resp.GetUsers() {
				fmt.Println(v)
			}

			cancel()
		},
	}

	searchCmd.Flags().String("address", "", "Address the server will bind to; alternatively, use ADDRESS var")
	_ = viper.BindPFlag("address", searchCmd.Flags().Lookup("address"))

	rootCmd.AddCommand(searchCmd)
}
