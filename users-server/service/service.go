package service

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/go-pg/pg/orm"

	"github.com/go-pg/pg"

	pb "github.com/maelvls/users-grpc/schema/user"
	"github.com/rs/xid"
	context "golang.org/x/net/context"
)

type User struct {
	Id        string `json:"id,omitempty"`
	Age       int32  `json:"age,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty" pg:"unique"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

func toPB(u User) pb.User {
	return pb.User{
		Id:      u.Id,
		Age:     u.Age,
		Name:    &pb.Name{First: u.FirstName, Last: u.FirstName},
		Email:   u.Email,
		Phone:   u.Phone,
		Address: u.Address,
	}
}

func toPBs(users []User) []pb.User {
	users2 := make([]pb.User, 0, len(users))
	for _, user := range users {
		users2 = append(users2, toPB(user))
	}
	return users2
}

func fromPB(u pb.User) User {
	return User{
		Id:        u.Id,
		Age:       u.Age,
		FirstName: u.Name.First,
		LastName:  u.Name.Last,
		Email:     u.Email,
		Phone:     u.Phone,
		Address:   u.Address,
	}
}

func fromPBs(users []pb.User) []User {
	users2 := make([]User, 0, len(users))
	for _, user := range users {
		users2 = append(users2, fromPB(user))
	}
	return users2
}

// MemDB is a simple in-memory DB by Hashicorp. As I wanted to keep things
// simple, I did not go with Postgres.

// InitSchema initializes the DB schema.
func InitSchema(db *pg.DB) error {
	// Create the DB schema for each Go type. We only have one right now, so
	// this loop might be slightly overkill.
	for _, model := range []interface{}{(*pb.User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// UserImpl implements my users-grpc service. No constructor needed.
type UserImpl struct {
}

// Create a user. If the given user has no id, generate one.
func (svc UserImpl) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateResp, error) {
	tx, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}

	user := req.GetUser()
	if user.Id == "" {
		user.Id = xid.New().String()
	}

	if err := tx.Insert(user); err != nil {
		panic(err)
	}

	resp := &pb.CreateResp{User: user, Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// List all users.
func (svc UserImpl) List(ctx context.Context, req *pb.ListReq) (*pb.SearchResp, error) {
	tx, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}

	var users []pb.User
	err = tx.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	resp := &pb.SearchResp{Users: pointerize(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchAge searches all users in the range [from, to_included].
func (svc UserImpl) SearchAge(ctx context.Context, req *pb.SearchAgeReq) (*pb.SearchResp, error) {

	if req.AgeRange == nil {
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  "field AgeRange{From: int, ToIncluded: int} missing"},
		}, nil
	}
	if req.AgeRange.From > req.AgeRange.ToIncluded {
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  "the From field must be lower or equal to ToIncluded"},
		}, nil
	}

	tx, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}

	// Select users in that age range.
	var users []pb.User
	_, err = tx.Query(&users, `SELECT * FROM users WHERE ? <= age AND age <= ?
    `, req.AgeRange.From, req.AgeRange.ToIncluded)
	if err != nil {
		panic(err)
	}

	resp := &pb.SearchResp{Users: pointerize(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchName searches a user by a part of its first or last name.
func (svc UserImpl) SearchName(ctx context.Context, req *pb.SearchNameReq) (*pb.SearchResp, error) {

	if req.Query == "" {
		return &pb.SearchResp{Users: make([]*pb.User, 0), Status: &pb.Status{
			Code: pb.Status_INVALID_QUERY,
			Msg:  "query cannot be empty"},
		}, nil
	}

	tx, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}

	// Select users in that age range.
	var users []pb.User
	_, err = tx.Query(&users, `
		SELECT * FROM users
		WHERE firstName LIKE %?%
		   OR lastName LIKE %?%
		`, req.Query, req.Query)

	if err != nil {
		return nil, fmt.Errorf("error when getting data from db: %e", err)
	}

	resp := &pb.SearchResp{Users: pointerize(users), Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// GetByEmail returns a user by its email.
func (svc UserImpl) GetByEmail(ctx context.Context, req *pb.GetByEmailReq) (*pb.GetByEmailResp, error) {

	user := &pb.User{Email: req.Email}
	tx, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}

	err = tx.Select(&user)
	if err != nil {
		panic(err)
	}

	// When not found, gracefully return 'email not found'.
	if user == nil {
		return &pb.GetByEmailResp{User: &pb.User{}, Status: &pb.Status{
			Code: pb.Status_FAILED,
			Msg:  "email not found"},
		}, nil
	}

	resp := &pb.GetByEmailResp{User: user, Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// Go from `[]pb.User` to `[]*pb.User`.
func pointerize(users []pb.User) []*pb.User {
	var users2 = make([]*pb.User, 0, len(users))
	for _, user := range users {
		// The iteration variables use the same variable address, which means
		// we can't directly use the pointer of 'u' here; we first need to
		// create a new local variable.
		// See: https://stackoverflow.com/questions/52980172
		u := user
		users2 = append(users2, &u)
	}
	return users2
}

type dummy int

// unexported: it is just a type + var for serving as a key for
// `context.WithValue`.
var txnKey dummy

// NewContext returns a new Context that carries value txn.
func NewContext(parent context.Context, txn *pg.Tx) context.Context {
	return context.WithValue(parent, txnKey, txn)
}

// FromContext returns the *Transaction value stored in ctx, if any.
func FromContext(ctx context.Context) (*pg.Tx, error) {
	tx, ok := ctx.Value(txnKey).(*pg.Tx)

	if !ok {
		return nil, fmt.Errorf(`transaction was not found in the gRPC context;
		it should be set by either grpc.TxInterceptor or manually (in tests)`)
	}

	return tx, nil
}

// TxInterceptor creates an interceptor that wraps the calls in DB
// transactions. I got this idea from gorm.UnaryServerInterceptor (gorm
// doesn't have that by default but infobloxopen/atlas-app-toolkit
// implements it). See
// https://github.com/infobloxopen/atlas-app-toolkit/blob/master/gorm/README.md#transaction-management
// and
// https://medium.com/@shijuvar/writing-grpc-interceptors-in-go-bf3e7671fe48
func TxInterceptor(db *pg.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tx, err := db.Begin()
		defer func() {
			_ = tx.Rollback()
		}()
		if err != nil {
			return nil, fmt.Errorf("could not begin transaction before handling the request: %e", err)
		}

		// Add the transaction handle to the gRPC context.
		ctx = NewContext(ctx, tx)

		// Continue the work
		h, err := handler(ctx, req)
		return h, err
	}
}
