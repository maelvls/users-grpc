package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/maelvls/users-grpc/schema/user"
	pb "github.com/maelvls/users-grpc/schema/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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
			cc, err := grpc.Dial(client.address, grpc.WithInsecure())
			if err != nil {
				logrus.Errorf("grpc client: %v\n", err)
				os.Exit(1)
			}

			client := user.NewUserServiceClient(cc)
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

			if err != nil {
				logrus.Errorf("grpc client: %v\n", err)
				os.Exit(1)
			}

			if resp.GetStatus().GetCode() != user.Status_SUCCESS {
				logrus.Errorf("grpc client: %v\n", resp.GetStatus())
				os.Exit(1)
			}

			// fmt.Println(Spprint(resp.GetUser()))
		},
	}

	createCmd.Flags().String("firstname", "", "") // Brianna
	createCmd.Flags().String("lastname", "", "")  // Shelton
	createCmd.Flags().String("email", "", "")     // brianna.shelton@email.org
	createCmd.Flags().Int32("age", 0, "")
	createCmd.Flags().String("postaladdress", "", "") // 255 Cortelyou Road, Volta, Indiana, 1608

	rootCmd.AddCommand(createCmd)
}
