package grpc

import (
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"

	service "github.com/maelvls/users-grpc/pkg/service"
	pb "github.com/maelvls/users-grpc/schema/user"
)

// For testing purposes.
type UserService interface {
	Create(*memdb.Txn, service.User) error
	List(*memdb.Txn) ([]service.User, error)
	SearchAge(txn *memdb.Txn, ageFrom, ageTo int32) ([]service.User, error)
	SearchName(txn *memdb.Txn, query string) ([]service.User, error)
	GetByEmail(txn *memdb.Txn, email string) (service.User, error)
}

// UserServer implements the GRPC endpoints of the "user" service. If I
// also wanted to be able to trace my service (e.g. using jaeger), I would
// also make sure to store opentracing.Tracer there.
type UserServer struct {
	DB  *memdb.MemDB
	Svc UserService // For testing purposes.
}

// NewUserServer returns a new server.
func NewUserServer() *UserServer {
	return &UserServer{DB: service.NewDBOrPanic(), Svc: service.UserSvc{}}
}

// Create a user. If the given user has no id, generate one.
func (server *UserServer) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateResp, error) {
	txn := server.DB.Txn(true)
	defer txn.Abort()

	err := server.Svc.Create(txn, FromPB(req.User))

	status := &pb.Status{Code: pb.Status_SUCCESS}
	switch {
	case err == service.EmailAlreadyExists:
		status = &pb.Status{Code: pb.Status_FAILED, Msg: err.Error()}
	case err != nil:
		logrus.WithError(err).WithField("email", req.User.Email).Error("Create returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while creating user, email=" + req.User.Email)
	}

	user, err := server.Svc.GetByEmail(txn, req.User.Email)
	if err != nil {
		logrus.Error("GetByEmail returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while finding the user, email=" + req.User.Email)
	}

	txn.Commit()
	return &pb.CreateResp{User: ToPB(user), Status: status}, nil
}

// List all users.
func (server *UserServer) List(ctx context.Context, req *pb.ListReq) (*pb.SearchResp, error) {
	txn := server.DB.Txn(false) // read-only transaction
	defer txn.Abort()

	users, err := server.Svc.List(txn)
	if err != nil {
		logrus.WithError(err).Error("List returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while listing")
	}

	resp := &pb.SearchResp{Users: ToPBs(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchAge searches all users in the range [from, to_included].
func (server *UserServer) SearchAge(ctx context.Context, req *pb.SearchAgeReq) (*pb.SearchResp, error) {
	if req.AgeRange == nil {
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  ""},
		}, nil
	}

	txn := server.DB.Txn(false)
	defer txn.Abort()

	users, err := server.Svc.SearchAge(txn, req.AgeRange.From, req.AgeRange.ToIncluded)

	switch {
	case err == service.AgeFromIsGreaterThanAgeTo:
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  "the From field must be lower or equal to ToIncluded"},
		}, nil
	}

	resp := &pb.SearchResp{Users: ToPBs(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchName searches a user by a part of its first or last name.
func (server *UserServer) SearchName(ctx context.Context, req *pb.SearchNameReq) (*pb.SearchResp, error) {
	txn := server.DB.Txn(false)
	defer txn.Abort()

	users, err := server.Svc.SearchName(txn, req.Query)
	switch {
	case err == service.NameQueryEmpty:
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  "name query cannot be empty"},
		}, nil
	case err != nil:
		logrus.WithError(err).WithField("query", req.Query).Error("SearchName returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while finding users by name, query=" + req.Query)
	}

	return &pb.SearchResp{Users: ToPBs(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}, nil
}

// GetByEmail returns a user by its email.
func (server *UserServer) GetByEmail(ctx context.Context, req *pb.GetByEmailReq) (*pb.GetByEmailResp, error) {
	txn := server.DB.Txn(false)
	defer txn.Abort()

	user, err := server.Svc.GetByEmail(txn, req.Email)
	switch {
	case err == service.EmailNotFound:
		return &pb.GetByEmailResp{User: &pb.User{}, Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  "this email cannot be found"},
		}, nil
	case err != nil:
		logrus.WithError(err).WithField("email", req.Email).Error("GetByEmail returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while getting a user by its email, email=" + req.Email)
	}

	resp := &pb.GetByEmailResp{User: ToPB(user), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

func FromPB(u *pb.User) service.User {
	return service.User{
		ID:        u.Id,
		Age:       u.Age,
		FirstName: u.Name.First,
		LastName:  u.Name.Last,
		Email:     u.Email,
		Phone:     u.Phone,
		Address:   u.Address,
	}
}

func ToPB(u service.User) *pb.User {
	return &pb.User{
		Id:      u.ID,
		Age:     u.Age,
		Name:    &pb.Name{First: u.FirstName, Last: u.LastName},
		Email:   u.Email,
		Phone:   u.Phone,
		Address: u.Address,
	}
}

func ToPBs(users []service.User) []*pb.User {
	users2 := make([]*pb.User, 0, len(users))
	for _, user := range users {
		users2 = append(users2, ToPB(user))
	}
	return users2
}
