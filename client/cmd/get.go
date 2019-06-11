package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/maelvls/quote/schema/user"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	getCmd := &cobra.Command{
		Use:   "get EMAIL",
		Short: "prints an user by its email (must be exact, not partial)",
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
				logrus.Fatalf("grpc client: %v\n", err)
			}

			client := user.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.GetByEmail(ctx, &user.GetByEmailReq{Email: givenEmail})

			if err != nil {
				logrus.Fatalf("grpc client: %v\n", err)
			}

			if resp.GetStatus().GetCode() == user.Status_FAILED {
				logrus.Fatalf("email not found")
			}

			if resp.GetStatus().GetCode() != user.Status_SUCCESS {
				logrus.Fatalf("grpc client: %#+v", resp.GetStatus())
			}

			fmt.Println(Spprint(resp.User))
		},
	}
	rootCmd.AddCommand(getCmd)
}
