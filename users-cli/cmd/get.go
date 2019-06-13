package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/maelvls/users-grpc/schema/user"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	getCmd := &cobra.Command{
		Use:   "get EMAIL",
		Short: "Fetch a user by its email (must be exact, not partial)",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires an email as argument")
			}
			return nil
		},
		Run: func(getCmd *cobra.Command, args []string) {
			givenEmail := args[0]

			cc, err := grpc.Dial(client.address, grpc.WithInsecure())
			if err != nil {
				logrus.Errorf("grpc client: %v\n", err)
				os.Exit(1)
			}

			client := user.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.GetByEmail(ctx, &user.GetByEmailReq{Email: givenEmail})

			if err != nil {
				logrus.Errorf("grpc client: %v\n", err)
				os.Exit(1)
			}

			if resp.GetStatus().GetCode() == user.Status_FAILED {
				logrus.Errorf("email not found")
				os.Exit(1)
			}

			if resp.GetStatus().GetCode() != user.Status_SUCCESS {
				logrus.Errorf("grpc client: %#+v", resp.GetStatus())
				os.Exit(1)
			}

			fmt.Println(Spprint(resp.User))
		},
	}
	rootCmd.AddCommand(getCmd)
}
