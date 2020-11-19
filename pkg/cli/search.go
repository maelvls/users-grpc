package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/maelvls/users-grpc/pkg/cli/logutil"
	pb "github.com/maelvls/users-grpc/schema/user"
	"github.com/spf13/cobra"
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search (--name=PARTIALNAME | --agefrom=N --ageto=M)",
		Short: "Search users from the remote users-server",
		Run: func(searchCmd *cobra.Command, args []string) {
			client, err := createClient(cfg)
			if err != nil {
				logutil.Errorf("%v", err)
				os.Exit(1)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var searchAge, searchName bool
			name, _ := searchCmd.Flags().GetString("name")
			if name != "" {
				searchName = true
			}

			ageFrom, err := searchCmd.Flags().GetInt32("agefrom")
			if err != nil {
				logutil.Errorf("--agefrom is not a number")
				os.Exit(1)
			}
			ageTo, err := searchCmd.Flags().GetInt32("ageto")
			if err != nil {
				logutil.Errorf("--ageto is not a number")
				os.Exit(1)
			}
			if ageTo != 0 || ageFrom != 0 {
				searchAge = true
			}

			switch {
			case searchAge && searchName:
				logutil.Errorf("cannot search by age and by name at the same time")
			case !searchAge && !searchName:
				logutil.Errorf("need one of '--name=PARTIALNAME' or '--agefrom=N' + '--ageto=M'")
			case searchAge:
				req := &pb.SearchAgeReq{AgeRange: &pb.SearchAgeReq_AgeRange{From: int32(ageFrom), ToIncluded: int32(ageTo)}}

				resp, err := client.SearchAge(ctx, req)
				switch {
				case err != nil:
					logutil.Errorf("searching age: %v", err)
					os.Exit(1)
				case resp.GetStatus().GetCode() != pb.Status_SUCCESS:
					logutil.Errorf("%s: %s", resp.Status.Code, resp.Status.Msg)
					os.Exit(1)
				default:
					// Happy path continuing below.
				}

				for _, u := range resp.GetUsers() {
					fmt.Println(Spprint(u))
				}
			case searchName:
				resp, err := client.SearchName(ctx, &pb.SearchNameReq{Query: name})
				switch {
				case err != nil:
					logutil.Errorf("searching by name: %v", err)
					os.Exit(1)
				case resp.GetStatus().GetCode() != pb.Status_SUCCESS:
					logutil.Errorf("%s: %s", resp.Status.Code, resp.Status.Msg)
					os.Exit(1)
				default:
					// Happy path continuing below.
				}

				for _, u := range resp.GetUsers() {
					fmt.Println(Spprint(u))
				}
			}
		},
	}

	searchCmd.Flags().String("name", "", "Search with a substring of first or last name; search is case-insensitive and special characters insensitive (e.g., searching 'mael' will return 'Maël') // brianna.shelton@email.org")
	searchCmd.Flags().Int32("agefrom", 0, "Search in [agefrom, ageto]")
	searchCmd.Flags().Int32("ageto", 0, "")

	rootCmd.AddCommand(searchCmd)
}
