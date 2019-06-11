package service

import (
	"context"
	"testing"

	memdb "github.com/hashicorp/go-memdb"
	pb "github.com/maelvls/users-grpc/schema/user"
	td "github.com/maxatome/go-testdeep"
)

// Helper for bootstraping a DB with the given users.
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
	// I don't really know how to test that
}

func TestNewUserImpl(t *testing.T) {
	svc := NewUserImpl()
	td.CmpStruct(t, svc, (*UserImpl)(nil), td.StructFields{"DB": td.NotNil()}, "DB isn't nil")
}

func TestUserImpl_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateReq
	}
	tests := []struct {
		name        string
		svc         *UserImpl
		args        args
		want        *pb.CreateResp
		fieldChecks td.StructFields
		postChecks  func(t *testing.T, svc *UserImpl)
	}{
		{
			name: "when a user is created, it should appear in the DB",
			svc: &UserImpl{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
			})},
			args: args{req: &pb.CreateReq{
				User: &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}},
			},
			want: &pb.CreateResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				User: &pb.User{
					Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr",
				},
			},
			fieldChecks: td.StructFields{},
			postChecks: func(t *testing.T, svc *UserImpl) {
				// Check that the user exists
				txn := svc.DB.Txn(false)
				defer txn.Abort()
				raw, err := txn.First("user", "id", "zikuwcus@awobik.kr")
				if td.CmpNoError(t, err) {
					td.CmpNotNil(t, raw)
				}
			},
		},
		{
			name: "when a user is created with the 'Id' field missing, the Id should be generated",
			svc:  &UserImpl{DB: initDBWith([]*pb.User{})},
			args: args{req: &pb.CreateReq{
				User: &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Email: "zikuwcus@awobik.kr"}},
			},
			want:        nil,
			fieldChecks: td.StructFields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.svc.Create(tt.args.ctx, tt.args.req)
			if td.CmpNoError(t, err) {
				td.CmpStruct(t, got, tt.want, tt.fieldChecks, tt.name)
			}
			if tt.postChecks != nil {
				tt.postChecks(t, tt.svc)
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
		name string
		svc  *UserImpl
		args args
		want *pb.SearchResp
	}{
		{
			name: "with no DB record, List should return an empty list of users",
			svc:  NewUserImpl(),
			args: args{req: &pb.ListReq{}},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users:  make([]*pb.User, 0),
			},
		},
		{
			name: "with 3 users in DB, List should return a list of 3 users",
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
		name string
		svc  *UserImpl
		args args
		want *pb.SearchResp
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
		name string
		svc  *UserImpl
		args args
		want *pb.SearchResp
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
	// Meh, that's not very useful to test that I guess... for now
}

func TestUserImpl_GetByEmail(t *testing.T) {
	type fields struct {
		DB *memdb.MemDB
	}
	type args struct {
		ctx context.Context
		req *pb.GetByEmailReq
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *pb.GetByEmailResp
	}{
		{
			name: "should return FAILED when no user has this email",
			fields: fields{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "le@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.GetByEmailReq{Email: "someemail@gmail.com"}},
			want: &pb.GetByEmailResp{
				Status: &pb.Status{Code: pb.Status_FAILED, Msg: "email not found"},
			},
		},
		{
			name: "should return Wayne when 'wayne.keller@rec.gb' is given",
			fields: fields{DB: initDBWith([]*pb.User{
				{Name: &pb.Name{First: "Elnora", Last: "Morales"}, Age: 21, Id: "ba3d530", Email: "eza@pod.ru"},
				{Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "wayne.keller@rec.gb"},
				{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			})},
			args: args{req: &pb.GetByEmailReq{Email: "wayne.keller@rec.gb"}},
			want: &pb.GetByEmailResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				User: &pb.User{
					Name: &pb.Name{First: "Wayne", Last: "Keller"}, Age: 42, Id: "c7dca0a", Email: "wayne.keller@rec.gb",
				},
			},
		},
	}
	for _, tt := range tests {
		svc := &UserImpl{
			DB: tt.fields.DB,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.GetByEmail(tt.args.ctx, tt.args.req)
			if td.CmpNoError(t, err) {
				td.CmpStruct(t, got, tt.want, td.StructFields{}, tt.name)
			}
		})

	}
}
