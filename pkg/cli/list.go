package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/maelvls/users-grpc/pkg/cli/logutil"
	"github.com/maelvls/users-grpc/schema/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(listCmd *cobra.Command, args []string) {
			cc, err := grpc.Dial(client.address, grpc.WithInsecure())
			if err != nil {
				logutil.Errorf("grpc client: %v", err)
				os.Exit(1)
			}

			client := user.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			resp, err := client.List(ctx, &user.ListReq{})

			if err != nil {
				logutil.Errorf("grpc client: %v", err)
				os.Exit(1)
			}

			if resp.GetStatus().GetCode() != user.Status_SUCCESS {
				logutil.Errorf("grpc client: %v", resp.GetStatus())
				os.Exit(1)
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
