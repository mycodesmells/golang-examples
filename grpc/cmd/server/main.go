package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "github.com/mycodesmells/golang-examples/grpc/proto/service"
)

func main() {
	addr := ":6000"
	clientAddr := fmt.Sprintf("localhost%s", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	go runGRPC(lis)
	runHTTP(clientAddr)
}

func runGRPC(lis net.Listener) {
	creds, err := credentials.NewServerTLSFromFile("cmd/server/server-cert.pem", "cmd/server/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	server := grpc.NewServer(
		grpc.Creds(creds),
		// grpc.UnaryInterceptor(AuthInterceptor),
	)
	pb.RegisterSimpleServerServer(server, NewServer())

	log.Printf("gRPC Listening on %s\n", lis.Addr().String())
	server.Serve(lis)
}

func runHTTP(clientAddr string) {
	runtime.HTTPError = CustomHTTPError

	addr := ":6001"
	creds, err := credentials.NewClientTLSFromFile("cmd/server/server-cert.pem", "")
	if err != nil {
		log.Fatalf("gateway cert load error: %s", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	mux := runtime.NewServeMux()
	if err := pb.RegisterSimpleServerHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

type server struct {
	users map[string]pb.User
}

func NewServer() server {
	return server{
		users: make(map[string]pb.User),
	}
}

func (s server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*empty.Empty, error) {
	log.Println("Creating user...")
	user := req.GetUser()

	if user.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	if user.Role == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "role cannot be empty")
	}

	s.users[user.Username] = *user

	log.Println("User created!")
	return &empty.Empty{}, nil
}

func (s server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	log.Println("Getting user!")

	if req.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	u, exists := s.users[req.Username]
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}

	log.Println("User found!")
	return &u, nil
}

func (s server) GreetUser(ctx context.Context, req *pb.GreetUserRequest) (*pb.GreetUserResponse, error) {
	log.Println("Greeting user...")
	if req.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}
	if req.Greeting == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "greeting cannot be empty")
	}

	user, err := s.GetUser(ctx, &pb.GetUserRequest{Username: req.Username})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find matching user")
	}

	return &pb.GreetUserResponse{
		Greeting: fmt.Sprintf("%s, %s! You are a great %s!", strings.Title(req.Greeting), user.Username, user.Role),
	}, nil
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	if len(meta["authorization"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	if meta["authorization"][0] != "valid-token" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}

type errorBody struct {
	Err string `json:"error,omitempty"`
}

func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	jErr := json.NewEncoder(w).Encode(errorBody{
		Err: grpc.ErrorDesc(err),
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
