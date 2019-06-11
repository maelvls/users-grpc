package service

import (
	"context"
	"reflect"
	"testing"

	memdb "github.com/hashicorp/go-memdb"
	pb "github.com/maelvls/quote/schema/user"
	td "github.com/maxatome/go-testdeep"
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
			name: "With no DB record, List should return an empty list of users",
			svc:  NewUserImpl(),
			args: args{req: &pb.ListReq{}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users:  make([]*pb.User, 0),
			},
		},
		{
			name: "With 3 users in DB, List should return a list of 3 users",
			svc: &UserImpl{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.ListReq{}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users: []*pb.User{
					{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
					{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
					{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.List(tt.args.ctx, tt.args.req)

			if td.CmpNoError(t, err) {
				td.CmpStruct(t, got, tt.want, td.StructFields{}, tt.name)
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
		{
			name: "should return INVALID_QUERY when From greater than ToIncluded",
			svc: &UserImpl{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.SearchAgeReq{AgeRange: &pb.SearchAgeReq_AgeRange{From: 21, ToIncluded: 10}}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_INVALID_QUERY, Msg: "the From field must be lower or equal to ToIncluded"},
				Users:  make([]*pb.User, 0),
			},
		},
		{
			name: "should return the single user of age 21",
			svc: &UserImpl{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.SearchAgeReq{AgeRange: &pb.SearchAgeReq_AgeRange{From: 21, ToIncluded: 21}}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users: []*pb.User{
					{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.SearchAge(tt.args.ctx, tt.args.req)
			if td.CmpNoError(t, err) {
				td.CmpStruct(t, got, tt.want, td.StructFields{}, tt.name)
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
		{
			name: "should INVALID_QUERY when the given query is empty",
			svc:  &UserImpl{},
			args: args{req: &pb.SearchNameReq{Query: ""}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_INVALID_QUERY, Msg: "query cannot be empty"},
				Users:  make([]*pb.User, 0),
			},
		},
		{
			name: "should return an empty list of users when nothing is found",
			svc: &UserImpl{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.SearchNameReq{Query: "something-that-cannot-be-found"}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users:  make([]*pb.User, 0),
			},
		},
		{
			name: "should return 'Elnora' when 'nor' is searched",
			svc: &UserImpl{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.SearchNameReq{Query: "nor"}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users: []*pb.User{
					{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.SearchName(tt.args.ctx, tt.args.req)
			if td.CmpNoError(t, err) {
				td.CmpStruct(t, got, tt.want, td.StructFields{}, tt.name)
			}
		})
	}
}

func Test_initDBWith(t *testing.T) {
	type args struct {
		users []*pb.User
	}
	tests := []struct {
		name string
		args args
		want *memdb.MemDB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initDBWith(tt.args.users); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initDBWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
