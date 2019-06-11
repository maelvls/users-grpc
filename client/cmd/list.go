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
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "lists all users",
		Run: func(listCmd *cobra.Command, args []string) {
			cc, err := grpc.Dial(client.address, grpc.WithInsecure())
			if err != nil {
				logrus.Fatalf("grpc client: %v\n", err)
			}

			client := user.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			resp, err := client.List(ctx, &user.ListReq{})

			if err != nil {
				logrus.Fatalf("grpc client: %v\n", err)
			}

			if resp.GetStatus().GetCode() != user.Status_SUCCESS {
				logrus.Fatalf("grpc client: %v\n", resp.GetStatus())
			}

			logrus.Debugf("number of users received: %v", len(resp.GetUsers()))

			// Finally, we can display the users.
			for _, user := range resp.GetUsers() {
				fmt.Println(Spprint(user))
			}

			cancel()
		},
	}

	rootCmd.AddCommand(listCmd)
}
