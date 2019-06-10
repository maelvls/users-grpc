package service

import (
	"context"
	"reflect"
	"testing"

	memdb "github.com/hashicorp/go-memdb"
	pb "github.com/maelvls/quote/schema/user"
)

func initDBWith(users []*pb.User) *memdb.MemDB {
	db := NewDB()
	txn := db.Txn(true)
	for _, user := range users {
		if err := txn.Insert("user", user); err != nil {
			panic(err)
		}
	}
	txn.Commit()
	return db
}

func TestNewDB(t *testing.T) {
	tests := []struct {
		name string
		want *memdb.MemDB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserImpl(t *testing.T) {
	tests := []struct {
		name string
		want *UserImpl
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserImpl(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserImpl_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateReq
	}
	tests := []struct {
		name    string
		svc     *UserImpl
		args    args
		want    *pb.CreateResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserImpl.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserImpl_List(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ListReq
	}

	tests := []struct {
		name    string
		svc     *UserImpl
		args    args
		want    *pb.SearchResp
		wantErr bool
	}{
		{
			name: "Empty DB should return no users",
			svc:  NewUserImpl(),
			args: args{req: &pb.ListReq{}},
			want: &pb.SearchResp{
				Users:  make([]*pb.User, 0),
				Status: &pb.Status{Code: pb.Status_SUCCESS},
			},
		},
		{
			name: "With 2 users, DB should return 2 users (needs Go 1.12+)",
			svc: &UserImpl{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.ListReq{}},
			want: &pb.SearchResp{
				Users: []*pb.User{
					{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
					{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
					{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
				},
				Status: &pb.Status{Code: pb.Status_SUCCESS},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.List(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserImpl.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserImpl_SearchAge(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.SearchAgeReq
	}
	tests := []struct {
		name    string
		svc     *UserImpl
		args    args
		want    *pb.SearchResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.SearchAge(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.SearchAge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserImpl.SearchAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserImpl_SearchName(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.SearchNameReq
	}
	tests := []struct {
		name    string
		svc     *UserImpl
		args    args
		want    *pb.SearchResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.SearchName(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.SearchName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserImpl.SearchName() = %v, want %v", got, tt.want)
			}
		})
	}
}
