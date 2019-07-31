package service

import (
	"context"
	"os"
	"testing"

	"github.com/go-pg/pg"
	"github.com/maelvls/users-grpc/schema/user"
	pb "github.com/maelvls/users-grpc/schema/user"
	td "github.com/maxatome/go-testdeep"
)

var db *pg.DB

func init() {
	// Get the database information. By default, if en vars are empty,
	// go-pg will use the defaults.
	pgOpts := &pg.Options{
		Addr: "",
		User: "postgres",
	}
	if os.Getenv("PG_ADDR") != "" {
		pgOpts.Addr = os.Getenv("PG_ADDR")
	}
	if os.Getenv("PG_USER") != "" {
		pgOpts.User = os.Getenv("PG_USER")
	}
	db := pg.Connect(pgOpts)
	InitSchema(db)
}

// Helper for bootstraping a DB with the given users.
func initDBWith(users []*pb.User) func(*pg.Tx) {
	return func(tx *pg.Tx) {
		for _, user := range users {
			if err := tx.Insert(user); err != nil {
				panic(err)
			}
		}
	}
}

func TestNewDB(t *testing.T) {
	// I don't really know how to test that
}

func TestUserImpl_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateReq
	}
	tests := []struct {
		name        string
		fieldChecks td.StructFields
		given       func(tx *pg.Tx)
		when        func(tx *pg.Tx) args
		then        func(t *testing.T, tx *pg.Tx, got *user.CreateResp)
	}{
		{
			given: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
			}),
			when: func(tx *pg.Tx) args {
				return args{ctx: NewContext(context.Background(), tx), req: &pb.CreateReq{
					User: &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}},
				}
			},
			then: func(t *testing.T, tx *pg.Tx, got *pb.CreateResp) {
				td.CmpStruct(t, got, &pb.CreateResp{
					Status: &pb.Status{Code: pb.Status_SUCCESS},
					User: &pb.User{
						Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr",
					},
				}, td.StructFields{},
					"when a user is created, it should appear in the DB")
				// Check that the user exists
				user := &pb.User{Email: "zikuwcus@awobik.kr"}
				err := tx.Select(user)
				if td.CmpNoError(t, err) {
					td.CmpNotNil(t, user)
				}
			},
		},
		{
			name:  "when a user is created with the 'Id' field missing, the Id should be generated",
			given: initDBWith([]*pb.User{}),
			when: func(tx *pg.Tx) args {
				return args{ctx: NewContext(context.Background(), tx), req: &pb.CreateReq{
					User: &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Email: "zikuwcus@awobik.kr"}},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := db.Begin()
			svc := &UserImpl{}
			tt.given(tx)
			args := tt.when(tx)
			got, err2 := svc.Create(args.ctx, args.req)
			if td.CmpNoError(t, err) && td.CmpNoError(t, err2) {
				tt.then(t, tx, got)
			}
			tx.Rollback()
		})
	}
}
