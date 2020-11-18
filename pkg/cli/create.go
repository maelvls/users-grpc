package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/maelvls/users-grpc/pkg/cli/logutil"

	"github.com/maelvls/users-grpc/schema/user"
	pb "github.com/maelvls/users-grpc/schema/user"
	"github.com/spf13/cobra"
)

func init() {
	createCmd := &cobra.Command{
		Use:   "create --email=EMAIL [--firstname] [--lastname] [--age] [--postaladdress]",
		Short: "Create a user",
		Args: func(createCmd *cobra.Command, args []string) error {
			email, err := createCmd.Flags().GetString("email")
			if email == "" || err != nil {
				return fmt.Errorf("--email=EMAIL required")
			}
			return nil
		},
		Run: func(createCmd *cobra.Command, args []string) {
			client, err := createClient(cfg)
			if err != nil {
				logutil.Errorf("grpc client: %v", err)
				os.Exit(1)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			firstname, _ := createCmd.Flags().GetString("firstname")
			lastname, _ := createCmd.Flags().GetString("lastname")
			age, _ := createCmd.Flags().GetInt32("age")

			postaladdress, _ := createCmd.Flags().GetString("postaladdress")
			email, _ := createCmd.Flags().GetString("email")

			usr := &pb.User{
				Email: email,
				Name: &pb.Name{
					First: firstname,
					Last:  lastname,
				},
				Age:     age,
				Address: postaladdress,
			}

			// Create the user.
			resp, err := client.Create(ctx, &user.CreateReq{User: usr})
			switch {
			case err != nil:
				logutil.Errorf("%v", err)
				os.Exit(1)
			case resp.GetStatus().GetCode() != user.Status_SUCCESS:
				logutil.Errorf("%v", resp.GetStatus())
				os.Exit(1)
			default:
				// Nothing.
			}
		},
	}

	createCmd.Flags().String("firstname", "", "") // Brianna
	createCmd.Flags().String("lastname", "", "")  // Shelton
	createCmd.Flags().String("email", "", "")     // brianna.shelton@email.org
	createCmd.Flags().Int32("age", 0, "")
	createCmd.Flags().String("postaladdress", "", "") // 255 Cortelyou Road, Volta, Indiana, 1608

	rootCmd.AddCommand(createCmd)
}
