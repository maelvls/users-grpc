package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/maelvls/quote/schema/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search (--name=PARTIALNAME | --agefrom=N --ageto=M)",
		Short: "searchs users from the remote quote service",
		Args: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			ageFromStr, _ := cmd.Flags().GetString("agefrom")
			ageToStr, _ := cmd.Flags().GetString("ageto")
			if name == "" && ageFromStr == "" && ageToStr == "" {
				return fmt.Errorf("need either '--name=PARTIALNAME' or '--agefrom=N' and '--ageto=M'")
			}
			if name != "" && (ageFromStr != "" || ageToStr != "") {
				return fmt.Errorf("cannot have both '--name=PARTIALNAME' or '--agefrom=N' and '--ageto=M'")
			}
			if (ageFromStr != "" && ageToStr == "") || (ageFromStr == "" && ageToStr != "") {
				return fmt.Errorf("cannot have both '--name=PARTIALNAME' or '--agefrom=N' and '--ageto=M'")
			}
			return nil
		},
		Run: func(searchCmd *cobra.Command, args []string) {

			cc, err := grpc.Dial(client.address, grpc.WithInsecure())
			if err != nil {
				logrus.Fatalf("grpc client: %v\n", err)
			}

			client := pb.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			name, _ := searchCmd.Flags().GetString("name")
			ageFromStr, _ := searchCmd.Flags().GetString("agefrom")
			ageToStr, err := searchCmd.Flags().GetString("ageto")

			// Check that either --agefrom --ageto OR --name is used, not
			// both at the same time.
			if ageFromStr != "" && ageToStr != "" && name != "" {
				logrus.Fatalf("cannot use both --agefrom/--ageto and --name")
			}

			if ageFromStr != "" && ageToStr != "" {
				var ageFrom int64
				var ageTo int64
				var err error

				// Parse both numbers
				if ageFromStr != "" {
					ageFrom, err = strconv.ParseInt(ageFromStr, 10, 32)
					if err != nil {
						logrus.Fatalf("--agefrom is not a number")
					}
				}
				if ageToStr != "" {
					ageTo, err = strconv.ParseInt(ageFromStr, 10, 32)
					if err != nil {
						logrus.Fatalf("--agefrom is not a number")
					}
				}
				resp, err := client.SearchAge(ctx, &pb.SearchAgeReq{
					AgeRange: &pb.SearchAgeReq_AgeRange{
						From:       int32(ageFrom),
						ToIncluded: int32(ageTo),
					},
				})

				if resp.GetStatus().GetCode() != pb.Status_SUCCESS {
					logrus.Fatalf("grpc client: %v\n", resp.GetStatus())
				}

				for _, u := range resp.GetUsers() {
					fmt.Println(Spprint(u))
				}
			}

			if name != "" {
				resp, err := client.SearchName(ctx, &pb.SearchNameReq{Query: name})
				if err != nil {
					logrus.Fatalf("grpc client: %v\n", err)
				}

				if resp.GetStatus().GetCode() != pb.Status_SUCCESS {
					logrus.Fatalf("grpc client: %v\n", resp.GetStatus())
				}

				for _, u := range resp.GetUsers() {
					fmt.Println(Spprint(u))
				}
			}

		},
	}

	searchCmd.Flags().String("name", "", "Search with a substring of first and last name") // brianna.shelton@undefined.org
	searchCmd.Flags().Int32("agefrom", 18, "Search in [agefrom, ageto]")
	searchCmd.Flags().Int32("ageto", 18, "")

	rootCmd.AddCommand(searchCmd)
}
