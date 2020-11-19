package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/maelvls/users-grpc/pkg/cli/logutil"
	"github.com/maelvls/users-grpc/schema/user"
	"github.com/spf13/cobra"
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

			client, err := createClient(cfg)
			if err != nil {
				logutil.Errorf("%v", err)
				os.Exit(1)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.GetByEmail(ctx, &user.GetByEmailReq{Email: givenEmail})
			switch {
			case err != nil:
				logutil.Errorf("get by email: %v", err)
				os.Exit(1)
			case resp.GetStatus().GetCode() == user.Status_FAILED:
				logutil.Errorf(resp.Status.Msg)
				os.Exit(1)
			case resp.GetStatus().GetCode() != user.Status_SUCCESS:
				logutil.Errorf("%#+v", resp.GetStatus())
				os.Exit(1)
			default:
				// Happy path continuing below.
			}

			// Finally, let's display the found user.
			fmt.Println(Spprint(resp.User))
		},
	}
	rootCmd.AddCommand(getCmd)
}
