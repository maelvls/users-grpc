package service

import (
	"fmt"
	"net"

	"github.com/google/uuid"
	pb_health "github.com/maelvls/quote/schema/health/v1"
	pb_quote "github.com/maelvls/quote/schema/quote"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type User struct {
	id       string
	username string
	name     string
}

var quotes map[uuid.UUID]string

// Server implements my quote service. If I also wanted to be able to trace my
// service (e.g. using jaeger), I would also make sure to store
// opentracing.Tracer there.
type Server struct {
}

// NewServer returns a new server.
func NewServer() *Server {
	return &Server{}
}

// Run starts the server
func (s *Server) Run(port int) error {
	srv := grpc.NewServer()
	pb_quote.RegisterQuoteServer(srv, s)
	pb_health.RegisterHealthServer(srv, &HealthImpl{})

	// Maybe we should let the user choose which address he wants to bind
	// to; in our case, when the host is unspecified (:80 is equivalent to
	// 0.0.0.0:80) then the local system. See: https://godoc.org/net#Dial
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return srv.Serve(lis)
}

// Search a quote using FTS (full-text search). For now, FTS is simply
// implemented using a regex but I plan on using actual natural language
// techniques for a more 'natural' querying (e.g., search synonyms).
func (s *Server) Search(ctx context.Context, req *pb_quote.SearchReq) (*pb_quote.SearchRes, error) {
	res := new(pb_quote.SearchRes)
	return res, nil
}

// Create a quote.
func (s *Server) Create(ctx context.Context, req *pb_quote.CreateReq) (*pb_quote.CreateRes, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	quotes[id] = req.GetQuote()

	res := new(pb_quote.CreateRes)
	return res, nil
}

var exampleData = map[string]interface{}{
	"a": "",
}

// HealthImpl implements the HealthServer interface (see health.pb.go file)
type HealthImpl struct{}

// Check returns SERVING.
func (h *HealthImpl) Check(ctx context.Context, args *pb_health.HealthCheckRequest) (*pb_health.HealthCheckResponse, error) {
	return &pb_health.HealthCheckResponse{
		Status: pb_health.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is not implemented for now.
func (h *HealthImpl) Watch(*pb_health.HealthCheckRequest, pb_health.Health_WatchServer) error {
	return nil
}
