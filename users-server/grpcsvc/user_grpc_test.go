package grpcsvc

import (
	"testing"

	pb "github.com/maelvls/users-grpc/schema/user"
	usersvc "github.com/maelvls/users-grpc/users-server/usersvc"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

func TestNewUserImpl(t *testing.T) {
	svc := NewUserImpl()
	td.CmpStruct(t, svc, (*UserImpl)(nil), td.StructFields{"DB": td.NotNil()}, "DB isn't nil")
}

func TestFromPB(t *testing.T) {
	given := pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}
	expect := usersvc.User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}

	assert.Equal(t, expect, FromPB(given))
}

func TestToPB(t *testing.T) {
	given := usersvc.User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}
	expect := &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}

	assert.Equal(t, expect, ToPB(given))
}