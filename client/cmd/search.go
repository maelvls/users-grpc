package cmd

import (
	"context"
	"fmt"
	"os"
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
		Run: func(searchCmd *cobra.Command, args []string) {

			cc, err := grpc.Dial(client.address, grpc.WithInsecure())
			if err != nil {
				logrus.Errorf("grpc client: %v\n", err)
				os.Exit(1)
			}

			client := pb.NewUserServiceClient(cc)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var mode = "SearchAge"

			name, _ := searchCmd.Flags().GetString("name")
			if name != "" {
				mode = "SearchName"
			}

			switch mode {
			case "":
				logrus.Errorf("need one of '--name=PARTIALNAME' or '--agefrom=N' + '--ageto=M'")
			case "SearchAge":
				var ageFrom, ageTo int32
				ageFrom, err = searchCmd.Flags().GetInt32("agefrom")
				if err != nil {
					logrus.Errorf("--agefrom is not a number")
					os.Exit(1)
				}
				ageTo, err = searchCmd.Flags().GetInt32("ageto")
				if err != nil {
					logrus.Errorf("--ageto is not a number")
					os.Exit(1)
				}
				resp, err := client.SearchAge(ctx, &pb.SearchAgeReq{
					AgeRange: &pb.SearchAgeReq_AgeRange{
						From:       int32(ageFrom),
						ToIncluded: int32(ageTo),
					},
				})

				if err != nil {
					logrus.Errorf("grpc client: %v\n", err)
					os.Exit(1)
				}

				if resp.GetStatus().GetCode() != pb.Status_SUCCESS {
					logrus.Errorf("grpc client: %v\n", resp.GetStatus())
					os.Exit(1)
				}

				for _, u := range resp.GetUsers() {
					fmt.Println(Spprint(u))
				}

			case "SearchName":
				resp, err := client.SearchName(ctx, &pb.SearchNameReq{Query: name})
				if err != nil {
					logrus.Errorf("grpc client: %v\n", err)
					os.Exit(1)
				}

				if resp.GetStatus().GetCode() != pb.Status_SUCCESS {
					logrus.Errorf("grpc client: %v\n", resp.GetStatus())
					os.Exit(1)
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
