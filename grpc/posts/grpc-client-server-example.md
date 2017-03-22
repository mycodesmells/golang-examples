# gRPC Client-Server Example

In the previous post, we started looking on gRPC by generating code using protocol buffers. As important as that step was, it would not be quite that powerful without some actual use case, where we'd like to have the same structure created on both sides od communication. This blog post provides you with such example use case - let's see how you can build your first gRPC server and its client in Golang.

### Communication model

In gRPC, we follow probably the most common communication pattern which is server-client, which is very similar to any HTTP communication. The flow starts with a server listening to the incoming requests on some TCP port. Then, whenever anything comes towards it this way, gets decoded to data structured generated from `*.proto` files. If the structure is parsed correctly, the server takes over with some request handling, then sends the result back in a binary form. If all goes well and both sides use the same proto definitions, the client should get the correct response.

### Generating service

We've already seen how to [generate code from proto](http://mycodesmells.com/post/intro-to-grpc---protocol-buffers), but when it comes to creating some gRPC services, we need something more. In order to generate code that can act as a server in our communication model, we need to define it using `service` and `rpc` keywords:

    ...
    import "google/protobuf/empty.proto";

    service SimpleServer {
        rpc CreateUser(CreateUserRequest) returns (google.protobuf.Empty) {}
        rpc GetUser(GetUserRequest) returns (User) {}
        rpc GreetUser(GreetUserRequest) returns (GreetUserResponse) {}
    }
    
This way we have a server named `SimpleServer` which has three functions that can be called remotely. The rest of the file needs to contain the definitions of all messages that are used as requests and responses. The lone exception in the example is `google.protobuf.Empty` which, you guessed it, is imported from an external package and is already defined.

**Note:** you may wonder why some responses are called `MethodNameResponse` and some are just an entity name. We've had a heated discussion at work about the way these objects should be named, we ended up following [GCP naming conventions](https://cloud.google.com/apis/design/naming_convention): return an entity for CRUD operations, return `MethodNameResponse` for everything else.

Once we define our service and the messages, we need to generate some code using `protoc`. This time we need to add a `plugin` to the `--go_out` option saying that we'd like to generate gRPC-specific code:

    $ protoc --go_out=plugins=grpc:${GOPATH}/src proto/service/service.proto

An output file from this command contains an interface that needs to be fulfilled by our implementation:

    type SimpleServerServer interface {
        CreateUser(context.Context, *CreateUserRequest) (*google_protobuf.Empty, error)
        GetUser(context.Context, *GetUserRequest) (*User, error)
        GreetUser(context.Context, *GreetUserRequest) (*GreetUserResponse, error)
    }
    
As you can see, each method has an extra argument of type `context.Context` which is very powerful, for example, it can be used to terminate the call after some timeout. Anyway, we can now create the implementation, Our `CreateUser` method could look like this:

    func (s server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*empty.Empty, error) {
        user := req.User
    
        if user.Username == "" {
            return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
        }
    
        if user.Role == "" {
            return nil, grpc.Errorf(codes.InvalidArgument, "role cannot be empty")
        }
    
        s.users[user.Username] = *user
        return &empty.Empty{}, nil
    }

As you can see, we create some user object and add it to the map, but this is just some implementation detail. As long as it implements an interface, it will be accessible from the gRPC client.

The last step is binding our server to the TCP connection:

    func main() {
        lis, err := net.Listen("tcp", "localhost:6000")
        if err != nil {
            log.Fatalf("failed to initializa TCP listen: %v", err)
        }
        defer lis.Close()
    
        server := grpc.NewServer()
        pb.RegisterSimpleServerServer(server, NewServer())
    
        server.Serve(lis)
    }

# Generating client

The other side of gRPC communication, the client, can be generated from proto as well. The output code from the `protoc` command listed above contains also something like:

    func NewSimpleServerClient(cc *grpc.ClientConn) SimpleServerClient {
        return &simpleServerClient{cc}
    }
    
So in order to make some gRPC calls we need to create a gRPC connection to some address and that's it! Let's to it, then:

    func main() {
        conn, err := grpc.Dial("localhost:6000", grpc.WithInsecure())
        if err != nil {
            log.Fatalf("Failed to start gRPC connection: %v", err)
        }
        defer conn.Close()
    
        client := pb.NewSimpleServerClient(conn)
        ...
    }
    
That's all you need to do to have a working client. Easy, right?

# Working example

In order to see how our server responds to the client's requests we need to start each application separately. First, let's see what requests we'll be making:


    // client/main.go
    func main() {
        ...
        client := pb.NewSimpleServerClient(conn)
    
        _, err = client.CreateUser(context.Background(), &pb.CreateUserRequest{User: &pb.User{Username: "slomek", Role: "joker"}})
        if err != nil {
            log.Fatalf("Failed to create user: %v", err)
        }
        log.Println("Created user!")
    
        resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Username: "slomek"})
        if err != nil {
            log.Fatalf("Failed to get created user: %v", err)
        }
        log.Printf("User exists: %v\n", resp)
    
        resp2, err := client.GreetUser(context.Background(), &pb.GreetUserRequest{Greeting: "howdy", Username: "slomek"})
        if err != nil {
            log.Fatalf("Failed to greet user: %v", err)
        }
        log.Printf("Greeting: %s\n", resp2.Greeting)
    }
    
So we'll create a user, then get it to make sure it has been saved properlya and finally we'll greet our new awesome user.

Let's start our server:

    $ go run cmd/server/main.go 
    2017/03/25 18:46:44 Listening on :6000

Once it's up, it's client's time to shine:

    $ go run cmd/client/main.go 
    2017/03/25 18:46:53 Created user!
    2017/03/25 18:46:53 User exists: username:"slomek" role:"joker" 
    2017/03/25 18:46:53 Greeting: Howdy, slomek! You are a great joker!

It works! Thank you gRPC!

The full working example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/grpc).
