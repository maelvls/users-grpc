package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/maelvls/quote/schema/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "searchs users from the remote quote service",
		Run: func(searchCmd *cobra.Command, args []string) {

			cc, err := grpc.Dial(client.address, grpc.WithInsecure())
			if err != nil {
				logrus.Fatalf("grpc client: %v\n", err)
			}

			client := user.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.SearchName(ctx, &user.SearchNameReq{Query: ""})

			if err != nil {
				logrus.Fatalf("grpc client: %v\n", err)
			}

			if resp.GetStatus().GetCode() != user.Status_SUCCESS {
				logrus.Fatalf("grpc client: %v\n", resp.GetStatus())
			}

			for _, u := range resp.GetUsers() {
				fmt.Println(Spprint(u))
			}

		},
	}

	searchCmd.Flags().String("name", "", "Search with a substring of first and last name") // brianna.shelton@undefined.org
	searchCmd.Flags().Int32("agefrom", 18, "Search in [agefrom, ageto]")
	searchCmd.Flags().Int32("ageto", 18, "")

	rootCmd.AddCommand(searchCmd)
}
