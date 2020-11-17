package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	memdb "github.com/hashicorp/go-memdb"
	"github.com/rs/xid"
)

var (
	EmailNotFound             = errors.New("email not found")
	EmailAlreadyExists        = errors.New("email already exists")
	NameQueryEmpty            = errors.New("name query cannot be empty")
	AgeFromIsGreaterThanAgeTo = errors.New("the starting age must be lower or equal to the ending age")
)

// MemDB is a simple in-memory DB by Hashicorp. As I wanted to keep things
// simple, I did not go with Postgres. I have a branch open with the
// Postgres implementation though:
// https://github.com/maelvls/users-grpc/pull/65

// NewDBOrPanic initializes the DB.
func NewDBOrPanic() *memdb.MemDB {
	// Create the DB schema.
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					// The primary key is on 'email'; we have to call this index 'id'
					// because go-memdb wants the table to have at least one 'id'
					// index.
					"id":    {Name: "id", Unique: true, Indexer: &memdb.StringFieldIndex{Field: "Email"}},
					"email": {Name: "email", Unique: true, Indexer: &memdb.StringFieldIndex{Field: "Email"}},
					"age":   {Name: "age", Unique: false, Indexer: &memdb.IntFieldIndex{Field: "Age"}},
				},
			},
		},
	}
	// Create a new data base.
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	return db
}

type User struct {
	ID        string `json:"id,omitempty"`
	Age       int32  `json:"age,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

// This struct is meant to make the service mockable for testing purposes.
// If I didn't need to test this, I would go with plain functions.
type UserSvc struct{}

// Create a user. The transaction must be created with write mode. If the
// given user has no ID, one will be generated randomly using the Mongo
// Object ID algorithm, see:
// https://docs.mongodb.com/manual/reference/method/ObjectId/
//
// The possible error is EmailAlreadyExists.
func (UserSvc) Create(txn *memdb.Txn, user User) error {
	if user.ID == "" {
		user.ID = xid.New().String()
	}

	// Let's make sure this email doesn't already exist.
	raw, err := txn.First("user", "email", user.Email)
	if err != nil {
		return fmt.Errorf("finding if the email %s is already used: %w", user.Email, err)
	}
	if raw != nil {
		return EmailAlreadyExists
	}

	err = txn.Insert("user", &user)
	if err != nil {
		return fmt.Errorf("inserting user %s: %w", user.Email, err)
	}

	return nil
}

// List all users.
func (UserSvc) List(txn *memdb.Txn) ([]User, error) {
	it, err := txn.Get("user", "email")
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	var users []User
	for raw := it.Next(); raw != nil; raw = it.Next() {
		user := raw.(*User)
		users = append(users, *user)
	}

	return users, nil
}

// SearchAge searches all users in the range [from, to_included]. The
// possible error is AgeFromIsGreaterThanAgeTo.
func (UserSvc) SearchAge(txn *memdb.Txn, ageFrom, ageTo int32) ([]User, error) {
	if ageFrom > ageTo {
		return nil, AgeFromIsGreaterThanAgeTo
	}

	// Range scan over people with ages in a range.
	it, err := txn.LowerBound("user", "age", ageFrom)
	if err != nil {
		return nil, fmt.Errorf("listing users starting at age %d: %w", ageFrom, err)
	}

	var users []User
	for raw := it.Next(); raw != nil; raw = it.Next() {
		u := raw.(*User)
		// Filter out all users that beyond the upper limit.
		if u.Age > ageTo {
			break
		}
		users = append(users, *u)
	}

	return users, nil
}

// SearchName searches a user by a part of its first or last name.
// The possible error is NameQueryEmpty.
func (UserSvc) SearchName(txn *memdb.Txn, query string) ([]User, error) {
	if query == "" {
		return nil, NameQueryEmpty
	}
	// This function filters out all users that do not contain the given
	// substr. Elmts are filtered/skipped when this function returns true.
	// This function should return false when an element should be kept.
	filterByFirstOrLastName := func(query string) func(interface{}) bool {
		return func(raw interface{}) bool {
			u, ok := raw.(*User)
			if !ok {
				logrus.Errorf("could not unpack a User, instead got: %#+v", raw)
				return true // Skip this element.
			}

			hasSubstr := strings.Contains(u.FirstName, query) ||
				strings.Contains(u.LastName, query)
			// We skip the element whenever the substr has not been matched.
			pleaseSkipIt := !hasSubstr
			return pleaseSkipIt
		}
	}

	result, err := txn.Get("user", "email")
	if err != nil {
		return nil, fmt.Errorf("err when getting data from db: %e", err)
	}

	it := memdb.NewFilterIterator(result, filterByFirstOrLastName(query))

	var users []User
	for raw := it.Next(); raw != nil; raw = it.Next() {
		u := raw.(*User)
		users = append(users, *u)
	}

	return users, nil
}

// GetByEmail returns a user by its email. May return EmailNotFound.
func (UserSvc) GetByEmail(txn *memdb.Txn, email string) (User, error) {
	raw, err := txn.First("user", "email", email)

	if err != nil {
		return User{}, fmt.Errorf("finding the user with email %s: %w", email, err)
	}

	// When not found, gracefully return 'email not found'
	if raw == nil {
		return User{}, EmailNotFound
	}

	user, ok := raw.(*User)
	if !ok {
		return User{}, fmt.Errorf("could not unpack a User, instead got: %#+v", raw)
	}

	return *user, nil
}
