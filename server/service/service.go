package service

import (
	memdb "github.com/hashicorp/go-memdb"
	"github.com/maelvls/quote/schema/user"
	pb "github.com/maelvls/quote/schema/user"
	"github.com/rs/xid"
	context "golang.org/x/net/context"
)

// DB is a simple K/V in-memory DB. As I want to keep things simple, it
// uses a map[string] and directly uses the User from the protobuf.

// // Create a write transaction
// txn := db.Txn(true)

// // Insert some people
// people := []*Person{
// 	&Person{"joe@aol.com", "Joe", 30},
// 	&Person{"lucy@aol.com", "Lucy", 35},
// 	&Person{"tariq@aol.com", "Tariq", 21},
// 	&Person{"dorothy@aol.com", "Dorothy", 53},
// }
// for _, p := range people {
// 	if err := txn.Insert("person", p); err != nil {
// 		panic(err)
// 	}
// }

// // Commit the transaction
// txn.Commit()

// // Create read-only transaction
// txn = db.Txn(false)
// defer txn.Abort()

// // Lookup by email
// raw, err := txn.First("person", "id", "joe@aol.com")
// if err != nil {
// 	panic(err)
// }

// // Say hi!
// fmt.Printf("Hello %s!\n", raw.(*Person).Name)

// // List all the people
// it, err := txn.Get("person", "id")
// if err != nil {
// 	panic(err)
// }

// fmt.Println("All the people:")
// for obj := it.Next(); obj != nil; obj = it.Next() {
// 	p := obj.(*Person)
// 	fmt.Printf("  %s\n", p.Name)
// }

// // Range scan over people with ages between 25 and 35 inclusive
// it, err = txn.LowerBound("person", "age", 25)
// if err != nil {
// 	panic(err)
// }

// fmt.Println("People aged 25 - 35:")
// for obj := it.Next(); obj != nil; obj = it.Next() {
// 	p := obj.(*Person)
// 	if p.Age > 35 {
// 		break
// 	}
// 	fmt.Printf("  %s is aged %d\n", p.Name, p.Age)
// }

// NewDB initializes the DB.
func NewDB() *memdb.MemDB {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {Name: "id", Unique: true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": {Name: "age", Unique: false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
				},
			},
		},
	}
	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	return db
}

// UserImpl implements my quote service. If I also wanted to be able to trace my
// service (e.g. using jaeger), I would also make sure to store
// opentracing.Tracer there.
type UserImpl struct {
	DB *memdb.MemDB
}

// NewUserImpl returns a new server.
func NewUserImpl() *UserImpl {
	return &UserImpl{DB: NewDB()}
}

// Create a user. If the given user has no id, generate one.
func (svc *UserImpl) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateResp, error) {
	user := req.GetUser()
	if user.Id == "" {
		user.Id = xid.New().String()
	}
	txn := svc.DB.Txn(true)
	if err := txn.Insert("user", user); err != nil {
		panic(err)
	}
	txn.Commit()

	resp := &pb.CreateResp{Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// List all users.
func (svc *UserImpl) List(ctx context.Context, req *pb.ListReq) (*pb.SearchResp, error) {
	// List all the people.
	txn := svc.DB.Txn(false) // read-only transaction
	defer txn.Abort()
	it, err := txn.Get("user", "id")
	if err != nil {
		panic(err)
	}

	var users = make([]*pb.User, 0)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		user := obj.(*user.User)
		users = append(users, user)
	}

	resp := &pb.SearchResp{Users: users, Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchAge searches all users in the range [from, to_included].
func (svc *UserImpl) SearchAge(ctx context.Context, req *pb.SearchAgeReq) (*pb.SearchResp, error) {
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

	txn := svc.DB.Txn(false) // read-only transaction
	// Range scan over people with ages in a range
	it, err := txn.LowerBound("user", "age", req.AgeRange.From)
	if err != nil {
		panic(err)
	}

	var users = make([]*pb.User, 0)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		u := obj.(*user.User)
		if u.Age > 35 {
			break
		}
		users = append(users, u)
	}

	resp := &pb.SearchResp{Users: users, Status: &pb.Status{Code: pb.Status_SUCCESS}}
	return resp, nil
}

// SearchName searches a user by a part of its first or last name.
func (svc *UserImpl) SearchName(ctx context.Context, req *pb.SearchNameReq) (*pb.SearchResp, error) {
	resp := &pb.SearchResp{Status: &pb.Status{Code: pb.Status_NO_IMPL_YET, Msg: "SearchName not implemented"}}
	return resp, nil
}
