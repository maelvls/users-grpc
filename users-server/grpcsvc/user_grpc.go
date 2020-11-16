package grpcsvc

import (
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"

	pb "github.com/maelvls/users-grpc/schema/user"
	usersvc "github.com/maelvls/users-grpc/users-server/usersvc"
)

// UserImpl implements my users-grpc service. If I also wanted to be able
// to trace my service (e.g. using jaeger), I would also make sure to store
// opentracing.Tracer there.
type UserImpl struct {
	DB *memdb.MemDB
}

// NewUserImpl returns a new server.
func NewUserImpl() *UserImpl {
	return &UserImpl{DB: usersvc.NewDBOrPanic()}
}

// Create a user. If the given user has no id, generate one.
func (svc *UserImpl) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateResp, error) {
	txn := svc.DB.Txn(true)
	defer txn.Abort()

	err := usersvc.Create(txn, FromPB(*req.User))

	status := &pb.Status{Code: pb.Status_SUCCESS}
	switch {
	case err == usersvc.EmailAlreadyExists:
		status = &pb.Status{Code: pb.Status_FAILED, Msg: err.Error()}
	case err != nil:
		logrus.WithError(err).WithField("email", req.User.Email).Error("Create returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while creating user, email=" + req.User.Email)
	}

	user, err := usersvc.GetByEmail(txn, req.User.Email)
	if err != nil {
		logrus.Error("GetByEmail returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while finding the user, email=" + req.User.Email)
	}

	txn.Commit()
	return &pb.CreateResp{User: ToPB(user), Status: status}, nil
}

// List all users.
func (svc *UserImpl) List(ctx context.Context, req *pb.ListReq) (*pb.SearchResp, error) {
	txn := svc.DB.Txn(false) // read-only transaction
	defer txn.Abort()

	users, err := usersvc.List(txn)
	if err != nil {
		logrus.WithError(err).Error("List returned an unexpected error")
		return nil, fmt.Errorf("something wrong happened while listing")
	}

	resp := &pb.SearchResp{Users: ToPBs(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchAge searches all users in the range [from, to_included].
func (svc *UserImpl) SearchAge(ctx context.Context, req *pb.SearchAgeReq) (*pb.SearchResp, error) {
	if req.AgeRange == nil {
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  ""},
		}, nil
	}

	txn := svc.DB.Txn(false)
	defer txn.Abort()

	users, err := usersvc.SearchAge(txn, req.AgeRange.From, req.AgeRange.ToIncluded)

	switch {
	case err == usersvc.AgeFromIsGreaterThanAgeTo:
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  "the From field must be lower or equal to ToIncluded"},
		}, nil
	}

	resp := &pb.SearchResp{Users: ToPBs(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchName searches a user by a part of its first or last name.
func (svc *UserImpl) SearchName(ctx context.Context, req *pb.SearchNameReq) (*pb.SearchResp, error) {
	txn := svc.DB.Txn(false)
	defer txn.Abort()

	users, err := usersvc.SearchName(txn, req.Query)
	switch {
	case err == usersvc.NameQueryEmpty:
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
func (svc *UserImpl) GetByEmail(ctx context.Context, req *pb.GetByEmailReq) (*pb.GetByEmailResp, error) {
	txn := svc.DB.Txn(false)
	defer txn.Abort()

	user, err := usersvc.GetByEmail(txn, req.Email)
	switch {
	case err == usersvc.EmailNotFound:
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

func FromPB(u pb.User) usersvc.User {
	return usersvc.User{
		ID:        u.Id,
		Age:       u.Age,
		FirstName: u.Name.First,
		LastName:  u.Name.Last,
		Email:     u.Email,
		Phone:     u.Phone,
		Address:   u.Address,
	}
}

func ToPB(u usersvc.User) *pb.User {
	return &pb.User{
		Id:      u.ID,
		Age:     u.Age,
		Name:    &pb.Name{First: u.FirstName, Last: u.LastName},
		Email:   u.Email,
		Phone:   u.Phone,
		Address: u.Address,
	}
}

func ToPBs(users []usersvc.User) []*pb.User {
	users2 := make([]*pb.User, 0, len(users))
	for _, user := range users {
		users2 = append(users2, ToPB(user))
	}
	return users2
}
