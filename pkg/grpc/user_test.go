package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/maelvls/users-grpc/pkg/grpc/mocks"
	service "github.com/maelvls/users-grpc/pkg/service"
	pb "github.com/maelvls/users-grpc/schema/user"
	td "github.com/maxatome/go-testdeep"
)

func TestUserServer_Create(t *testing.T) {
	tests := []struct {
		name      string
		givenReq  *pb.CreateReq
		givenMock func(rec *mocks.MockUserServiceMockRecorder)
		want      *pb.CreateResp
		wantErr   error
	}{
		{
			name: "when a user is created, it creates it",
			givenReq: &pb.CreateReq{
				User: &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.
					Create(someTxn(), service.User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}).
					Return(nil)
				rec.
					GetByEmail(someTxn(), "zikuwcus@awobik.kr").
					Return(service.User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}, nil)

			},
			want: &pb.CreateResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				User:   &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			},
			wantErr: nil,
		},
		{
			name:     "when the user already exists, return an understandable message",
			givenReq: &pb.CreateReq{User: &pb.User{Name: &pb.Name{}, Email: "zikuwcus@awobik.kr"}},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.
					Create(someTxn(), service.User{Email: "zikuwcus@awobik.kr"}).
					Return(service.EmailAlreadyExists)
			},
			want:    &pb.CreateResp{User: &pb.User{}, Status: &pb.Status{Code: pb.Status_FAILED, Msg: "email already exists"}},
			wantErr: nil,
		},
		{
			name:     "unknown Create errors should error the grpc request and hide the actual err message",
			givenReq: &pb.CreateReq{User: &pb.User{Name: &pb.Name{}, Email: "foo@bar.io"}},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.Create(someTxn(), service.User{Email: "foo@bar.io"}).Return(fmt.Errorf("some random error"))
			},
			want:    nil,
			wantErr: fmt.Errorf("something wrong happened while creating user, email=foo@bar.io"),
		},
		{
			name:     "unknown GetByEmail errors should error the grpc request and hide the actual err message",
			givenReq: &pb.CreateReq{User: &pb.User{Name: &pb.Name{}, Email: "foo@bar.io"}},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.Create(someTxn(), service.User{Email: "foo@bar.io"}).Return(nil)
				rec.GetByEmail(someTxn(), "foo@bar.io").Return(service.User{}, fmt.Errorf("unknown error"))
			},
			want:    nil,
			wantErr: fmt.Errorf("something wrong happened while finding the user, email=foo@bar.io"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockUserSvc := mocks.NewMockUserService(ctl)
			tt.givenMock(mockUserSvc.EXPECT())

			svc := &UserServer{
				Txn:      func(b bool) *memdb.Txn { return nil },
				Commit:   func(m *memdb.Txn) {},
				Rollback: func(m *memdb.Txn) {},
				Svc:      mockUserSvc,
			}

			got, gotErr := svc.Create(context.Background(), tt.givenReq)

			if tt.wantErr != nil {
				td.Cmp(t, gotErr, tt.wantErr)
				return
			}
			if td.CmpNoError(t, gotErr) {
				td.Cmp(t, got, tt.want)
			}
		})
	}
}

func TestUserServer_List(t *testing.T) {
	tests := []struct {
		name      string
		givenReq  *pb.ListReq
		givenMock func(rec *mocks.MockUserServiceMockRecorder)
		want      *pb.SearchResp
		wantErr   error
	}{
		{
			name:     "should return the list of users",
			givenReq: &pb.ListReq{},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.
					List(someTxn()).
					Return([]service.User{{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}}, nil)
			},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users:  []*pb.User{{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}},
			},
			wantErr: nil,
		},
		{
			name:     "unknown listing errors should error the grpc request and hide the actual err message",
			givenReq: &pb.ListReq{},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.
					List(someTxn()).
					Return(nil, fmt.Errorf("unknown list error"))
			},
			want: &pb.SearchResp{
				Status: &pb.Status{Code: pb.Status_SUCCESS},
				Users:  []*pb.User{{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}},
			},
			wantErr: fmt.Errorf("something wrong happened while listing users"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockUserSvc := mocks.NewMockUserService(ctl)
			tt.givenMock(mockUserSvc.EXPECT())

			svc := &UserServer{
				Txn:      func(b bool) *memdb.Txn { return nil },
				Commit:   func(m *memdb.Txn) {},
				Rollback: func(m *memdb.Txn) {},
				Svc:      mockUserSvc,
			}

			got, gotErr := svc.List(context.Background(), tt.givenReq)

			if tt.wantErr != nil {
				td.Cmp(t, gotErr, tt.wantErr)
				return
			}
			if td.CmpNoError(t, gotErr) {
				td.Cmp(t, got, tt.want)
			}
		})
	}
}

func TestUserServer_SearchAge(t *testing.T) {
	tests := []struct {
		name      string
		givenReq  *pb.SearchAgeReq
		givenMock func(rec *mocks.MockUserServiceMockRecorder)
		want      *pb.SearchResp
		wantErr   error
	}{
		{
			name:     "returns any found users",
			givenReq: &pb.SearchAgeReq{AgeRange: &pb.SearchAgeReq_AgeRange{From: 35, ToIncluded: 38}},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.SearchAge(someTxn(), int32(35), int32(38)).Return([]service.User{{Age: 38, Email: "zikuwcus@awobik.kr"}}, nil)
			},
			want: &pb.SearchResp{Status: &pb.Status{Code: pb.Status_SUCCESS}, Users: []*pb.User{{Name: &pb.Name{}, Age: 38, Email: "zikuwcus@awobik.kr"}}},
		},
		{
			name:      "should return an understandable message when age are wrong",
			givenReq:  &pb.SearchAgeReq{AgeRange: nil},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {},
			want:      &pb.SearchResp{Status: &pb.Status{Code: pb.Status_INVALID_QUERY, Msg: "the AgeRange object cannot be omitted"}, Users: []*pb.User{}},
		},
		{
			name:     "should return an understandable message when AgeRange is omitted",
			givenReq: &pb.SearchAgeReq{AgeRange: &pb.SearchAgeReq_AgeRange{From: 38, ToIncluded: 30}},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.SearchAge(someTxn(), int32(38), int32(30)).Return(nil, service.AgeFromIsGreaterThanAgeTo)
			},
			want: &pb.SearchResp{Status: &pb.Status{Code: pb.Status_INVALID_QUERY, Msg: "age is invalid, the 'from' age must be lower or equal to the 'to' age"}, Users: []*pb.User{}},
		},
		{
			name:     "unknown errors should error the grpc request and hide the actual err message",
			givenReq: &pb.SearchAgeReq{AgeRange: &pb.SearchAgeReq_AgeRange{}},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.SearchAge(someTxn(), int32(0), int32(0)).Return(nil, fmt.Errorf("unknown error"))
			},
			want:    nil,
			wantErr: fmt.Errorf("something wrong happened while searching users with their age"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockUserSvc := mocks.NewMockUserService(ctl)
			tt.givenMock(mockUserSvc.EXPECT())

			svc := &UserServer{
				Txn:      func(b bool) *memdb.Txn { return nil },
				Commit:   func(m *memdb.Txn) {},
				Rollback: func(m *memdb.Txn) {},
				Svc:      mockUserSvc,
			}

			got, gotErr := svc.SearchAge(context.Background(), tt.givenReq)

			if tt.wantErr != nil {
				td.Cmp(t, gotErr, tt.wantErr)
				return
			}
			if td.CmpNoError(t, gotErr) {
				td.Cmp(t, got, tt.want)
			}
		})
	}
}

func TestNewUserServer(t *testing.T) {
	svc := NewUserServer()
	td.CmpStruct(t, svc, (*UserServer)(nil), td.StructFields{
		"Commit":   td.NotNil(),
		"Rollback": td.NotNil(),
		"Txn":      td.NotNil(),
		"Svc":      td.NotNil(),
	})
}

func TestFromPB(t *testing.T) {
	given := &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}
	expect := service.User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}

	td.Cmp(t, FromPB(given), expect)
}

func TestToPB(t *testing.T) {
	given := service.User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}
	expect := &pb.User{Name: &pb.Name{First: "Flora", Last: "Hale"}, Age: 38, Id: "a4bcd38", Email: "zikuwcus@awobik.kr"}

	td.Cmp(t, ToPB(given), expect)
}

func someTxn() gomock.Matcher {
	return gomock.Any()
}

func TestUserServer_SearchName(t *testing.T) {
	tests := []struct {
		name      string
		givenReq  *pb.SearchNameReq
		givenMock func(rec *mocks.MockUserServiceMockRecorder)
		want      *pb.SearchResp
		wantErr   error
	}{
		{
			name:     "returns any found users",
			givenReq: &pb.SearchNameReq{Query: "oba"},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.SearchName(someTxn(), "oba").Return([]service.User{{FirstName: "Foobar"}}, nil)
			},
			want: &pb.SearchResp{Status: &pb.Status{Code: pb.Status_SUCCESS}, Users: []*pb.User{{Name: &pb.Name{First: "Foobar"}}}},
		},
		{
			name:     "should return an understandable message when quert is empty",
			givenReq: &pb.SearchNameReq{Query: ""},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.SearchName(someTxn(), "").Return(nil, service.NameQueryEmpty)
			},
			want: &pb.SearchResp{Status: &pb.Status{Code: pb.Status_INVALID_QUERY, Msg: "name query cannot be empty"}, Users: []*pb.User{}},
		},
		{
			name:     "unknown errors should error the grpc request and hide the actual err message",
			givenReq: &pb.SearchNameReq{Query: "blah"},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.SearchName(someTxn(), "blah").Return(nil, fmt.Errorf("unknown error"))
			},
			want:    nil,
			wantErr: fmt.Errorf("something wrong happened while finding users by name, query=blah"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockUserSvc := mocks.NewMockUserService(ctl)
			tt.givenMock(mockUserSvc.EXPECT())

			svc := &UserServer{
				Txn:      func(b bool) *memdb.Txn { return nil },
				Commit:   func(m *memdb.Txn) {},
				Rollback: func(m *memdb.Txn) {},
				Svc:      mockUserSvc,
			}

			got, gotErr := svc.SearchName(context.Background(), tt.givenReq)

			if tt.wantErr != nil {
				td.Cmp(t, gotErr, tt.wantErr)
				return
			}
			if td.CmpNoError(t, gotErr) {
				td.Cmp(t, got, tt.want)
			}
		})
	}
}

func TestUserServer_GetByEmail(t *testing.T) {
	tests := []struct {
		name      string
		givenReq  *pb.GetByEmailReq
		givenMock func(rec *mocks.MockUserServiceMockRecorder)
		want      *pb.GetByEmailResp
		wantErr   error
	}{
		{
			name:     "returns a user",
			givenReq: &pb.GetByEmailReq{Email: "zikuwcus@awobik.kr"},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.GetByEmail(someTxn(), "zikuwcus@awobik.kr").Return(service.User{Email: "zikuwcus@awobik.kr"}, nil)
			},
			want: &pb.GetByEmailResp{Status: &pb.Status{Code: pb.Status_SUCCESS}, User: &pb.User{Email: "zikuwcus@awobik.kr", Name: &pb.Name{}}},
		},
		{
			name:     "should return an understandable message when this email does not exist",
			givenReq: &pb.GetByEmailReq{Email: "zikuwcus@awobik.kr"},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.GetByEmail(someTxn(), "zikuwcus@awobik.kr").Return(service.User{}, service.EmailNotFound)
			},
			want: &pb.GetByEmailResp{Status: &pb.Status{Code: pb.Status_INVALID_QUERY, Msg: "the email zikuwcus@awobik.kr cannot be found"}, User: &pb.User{}},
		},
		{
			name:     "unknown errors should error the grpc request and hide the actual err message",
			givenReq: &pb.GetByEmailReq{Email: "foo@bar.io"},
			givenMock: func(rec *mocks.MockUserServiceMockRecorder) {
				rec.GetByEmail(someTxn(), "foo@bar.io").Return(service.User{}, fmt.Errorf("unknown error"))
			},
			want:    nil,
			wantErr: fmt.Errorf("something wrong happened while getting a user by its email, email=foo@bar.io"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockUserSvc := mocks.NewMockUserService(ctl)
			tt.givenMock(mockUserSvc.EXPECT())

			svc := &UserServer{
				Txn:      func(b bool) *memdb.Txn { return nil },
				Commit:   func(m *memdb.Txn) {},
				Rollback: func(m *memdb.Txn) {},
				Svc:      mockUserSvc,
			}

			got, gotErr := svc.GetByEmail(context.Background(), tt.givenReq)

			if tt.wantErr != nil {
				td.Cmp(t, gotErr, tt.wantErr)
				return
			}
			if td.CmpNoError(t, gotErr) {
				td.Cmp(t, got, tt.want)
			}
		})
	}
}
