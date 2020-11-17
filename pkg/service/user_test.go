package service

import (
	"fmt"
	"testing"

	memdb "github.com/hashicorp/go-memdb"
	td "github.com/maxatome/go-testdeep/td"
)

// Helper for filling the DB with the given users. Transaction must be
// opened in write mode.
func fillDBWith(users []User) func(*memdb.Txn) {
	return func(txn *memdb.Txn) {
		for _, user := range users {
			u := user // Due to the fact that all iterations share the same pointer.
			if err := txn.Insert("user", &u); err != nil {
				panic(err)
			}
		}
	}
}

func TestNewDB(t *testing.T) {
	// I don't really know how to test that
}

func TestCreate(t *testing.T) {
	db := NewDBOrPanic()

	tests := []struct {
		name        string
		init        func(txn *memdb.Txn)
		createUser  User
		wantErr     error
		fieldChecks td.StructFields
		postChecks  func(t *testing.T, txn *memdb.Txn)
	}{
		{
			name: "when a user is created, it should appear in the DB",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
			}),
			createUser:  User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			fieldChecks: td.StructFields{},
			postChecks: func(t *testing.T, txn *memdb.Txn) {
				// Check that the user exists.
				raw, err := txn.First("user", "id", "zikuwcus@awobik.kr")
				if td.CmpNoError(t, err) && td.CmpNotNil(t, raw) {
					user := raw.(*User)
					td.Cmp(t, User{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"}, *user)
				}
			},
		},
		{
			name:        "when a user is created with the 'Id' field missing, the Id should be generated",
			init:        fillDBWith([]User{}),
			createUser:  User{FirstName: "Flora", LastName: "Hale", Age: 38, Email: "zikuwcus@awobik.kr"},
			wantErr:     nil,
			fieldChecks: td.StructFields{},
		},
		{
			name: "when a user is created with an email that already exists, it should fail",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
			}),
			createUser:  User{FirstName: "Elnora", LastName: "Morales", Age: 38, Email: "eza@pod.ru"},
			wantErr:     EmailAlreadyExists,
			fieldChecks: td.StructFields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := db.Txn(true)
			defer txn.Abort()

			tt.init(txn)

			gotErr := UserSvc{}.Create(txn, tt.createUser)

			if tt.wantErr != nil {
				td.Cmp(t, gotErr, tt.wantErr)
				return
			}
			td.CmpNoError(t, gotErr)
			if tt.postChecks != nil {
				tt.postChecks(t, txn)
			}
		})
	}
}

func TestList(t *testing.T) {
	db := NewDBOrPanic()

	tests := []struct {
		name    string
		init    func(txn *memdb.Txn)
		want    []User
		wantErr error
	}{
		{
			name: "with no DB record, List should return an empty list of users",
			init: fillDBWith(nil),
			want: nil,
		},
		{
			name: "with 3 users in DB, List should return a list of 3 users",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
				{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			}),
			want: []User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
				{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := db.Txn(true)
			defer txn.Abort()

			tt.init(txn)

			got, gotErr := UserSvc{}.List(txn)

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

func TestSearchAge(t *testing.T) {
	db := NewDBOrPanic()

	tests := []struct {
		name    string
		init    func(txn *memdb.Txn)
		ageFrom int32
		ageTo   int32
		want    []User
		wantErr error
	}{
		{
			name:    "should return and error when fromAge is above toAge",
			init:    fillDBWith(nil),
			ageFrom: 21,
			ageTo:   10,
			wantErr: fmt.Errorf("the starting age must be lower or equal to the ending age"),
		},
		{
			name: "should return the single user of age 21",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
				{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			}),
			ageFrom: 21,
			ageTo:   21,
			want: []User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := db.Txn(true)
			defer txn.Abort()

			tt.init(txn)

			got, gotErr := UserSvc{}.SearchAge(txn, tt.ageFrom, tt.ageTo)
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

func TestSearchName(t *testing.T) {
	db := NewDBOrPanic()

	tests := []struct {
		name       string
		init       func(txn *memdb.Txn)
		searchName string
		want       []User
		wantErr    error
	}{
		{
			name:       "should return error when the given query is empty",
			init:       fillDBWith(nil),
			searchName: "",
			wantErr:    fmt.Errorf("name query cannot be empty"),
		},
		{
			name: "should return an empty list of users when nothing is found",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
				{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			}),
			searchName: "something-that-cannot-be-found",
			want:       nil,
		},
		{
			name: "should return 'Elnora' when 'nor' is searched",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
				{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			}),
			searchName: "nor",
			want:       []User{{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := db.Txn(true)
			defer txn.Abort()

			tt.init(txn)

			got, gotErr := UserSvc{}.SearchName(txn, tt.searchName)
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

func Test_fillDBWith(t *testing.T) {
	// Meh, that's not very useful to test that I guess... for now
}

func TestGetByEmail(t *testing.T) {
	db := NewDBOrPanic()

	tests := []struct {
		name     string
		init     func(txn *memdb.Txn)
		getEmail string
		want     User
		wantErr  error
	}{
		{
			name: "should return an error when no user has this email",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "le@rec.gb"},
				{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			}),
			getEmail: "someemail@gmail.com",
			wantErr:  fmt.Errorf("email not found"),
		},
		{
			name: "should return Wayne when 'wayne.keller@rec.gb' is given",
			init: fillDBWith([]User{
				{FirstName: "Elnora", LastName: "Morales", Age: 21, ID: "ba3d530", Email: "eza@pod.ru"},
				{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "wayne.keller@rec.gb"},
				{FirstName: "Flora", LastName: "Hale", Age: 38, ID: "a4bcd38", Email: "zikuwcus@awobik.kr"},
			}),
			getEmail: "wayne.keller@rec.gb",
			want:     User{FirstName: "Wayne", LastName: "Keller", Age: 42, ID: "c7dca0a", Email: "wayne.keller@rec.gb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := db.Txn(true)
			defer txn.Abort()

			tt.init(txn)

			got, gotErr := UserSvc{}.GetByEmail(txn, tt.getEmail)
			if tt.wantErr != nil {
				td.Cmp(t, gotErr, tt.wantErr)
				return
			}
			if td.CmpNoError(t, gotErr) {
				td.CmpStruct(t, got, tt.want, td.StructFields{}, tt.name)
			}
		})
	}
}
