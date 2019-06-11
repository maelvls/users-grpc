package cmd

import (
	"fmt"

	pb "github.com/maelvls/users-grpc/schema/user"
	"github.com/mgutz/ansi"
)

// Spprint is a helper functions for nicely displaying users.
func Spprint(u *pb.User) string {
	yel := ansi.ColorFunc("yellow+b")
	gre := ansi.ColorFunc("green")
	ansi.Color(u.Name.First, ansi.Yellow)
	return fmt.Sprintf("%s %s <%s> (%v years old, address: %s)",
		yel(u.Name.First),
		yel(u.Name.Last),
		gre(u.Email),
		u.Age,
		u.Address)
}
