# HTTP server from gRPC with Gateway

Building a gRPC service is not that complicated, even if we want to add some security, custom interceptors and complex request or response messages. Unfortunately, it's still not as popular as we'd like, and most of the time we need to support HTTP as well. Building two identical APIs seems like a terrible idea, right? With gRPC Gateway you can save lots of time, so let's jump in.

### gRPC vs HTTP

The main advantage of gRPC is how fast it is, and the fact that the code is auto-generated from `.proto` files. On the other hand, HTTP provides us with simplicity and popularity - after all almost every developer has some experience with such APIs, it can be easily tested with Postman, etc. What if we could combine the two?

[gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway) allows you to expose gRPC API via HTTP: it generates a simple server that connects to our existing gRPC server, translates JSON payloads into messages and everything is handled by gRPC from there.

### Before we start

The first steps we need to take is installing dependencies that are required by the gateway to work:

    go get -u github.com/grpc-ecosystem/grpc-gateway/...

The other thing we need to change is the authorization header. In the previous example we've used a metadata field called _token_, but since we are going to create HTTP server, we need to put the token into _authorization_ header, and since the metadata is auto-translated to headers, we might as well use the same name here as well.

    func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        ...
        if len(meta["authorization"]) != 1 {
        ...
        if meta["authorization"][0] != "valid-token" {
        ...
        return handler(ctx, req)
    }

### Changes in .proto

You will be amazed how little you need to change to the proto files. Basically, we just need to add `option` to the `rpc` methods, that contains HTTP method to be used, the path on which the endpoint will be accessible, and the way the data is being passed (`body` or inside path). Our three endpoints should look like this:

    service SimpleServer {
        rpc CreateUser(CreateUserRequest) returns (google.protobuf.Empty) {
            option (google.api.http) = {
                post: "/users",
                body: "*"
            };
        }
        rpc GetUser(GetUserRequest) returns (User) {
            option (google.api.http) = {
                get: "/users/{username}"
            };
        }
        rpc GreetUser(GreetUserRequest) returns (GreetUserResponse) {
            option (google.api.http) = {
                post: "/users/{username}/greet"
                body: "*"
            };
        }
    }

As you can see, the first endpoint expects all data to be passed as a body in POST request, the second one takes only a single parameter from the path. The last one expects `username` to be present in the path, the rest of fields needs to be sent inside the body.

**Note:** Our generated HTTP server will ignore `username` passed inside the body in `GreetUser`, only the one from the path does matter.

The second change that needs to be made is the way we generate code from proto files - we just need to add another option to the `protoc` command - `--grpc-gateway_out=logtostderr=true:.` (that way we have some error output shown to the console, just in case).

Now running our task from Makefile results in two files being generated:

    $ ls proto/service/
    service.proto
    $ make gen_proto
    protoc -I. -I /Users/slomek/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:/Users/slomek/go/src proto/service/service.proto --grpc-gateway_out=logtostderr=true:.
    protoc --go_out=/Users/slomek/go/src proto/message/message.proto --grpc-gateway_out=logtostderr=true:.
    $ ls proto/service/
    service.pb.go        service.pb.gw.go    service.proto

### Usage

In order to expose an HTTP server, we need just a couple of lines thanks to the code imported from generated files:

    func runHTTP(clientAddr string) {
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

We basically need to create a `ServerMux` from `grpc/runtime`, then pass it to the magical function that binds it with gRPC address (the one we'd send to our clients) and exposes itself on some given port. Then we can incorporate it inside our `main.go`:

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

**Note:** Since both `grpc.Server.Serve(..)` and `http.ListenAndServe()` are blocking (they stop the execution until the server is working), we need to run one of them in a separate goroutine.

Now when we start the server, we can use it with curl:

    $ go run cmd/server/main.go 
    2017/04/07 18:45:19 HTTP Listening on :6001
    2017/04/07 18:45:19 gRPC Listening on [::]:6000

Creating user:

    $ curl --request POST \
    >   --url http://localhost:6001/users \
    >   --header 'authorization: valid-token' \
    >   --header 'content-type: application/json' \
    >   --data '{"user":{"username":"budmore","role":"meetup"}}'

    {}

Getting user:

    $ curl --request GET \
    >   --url http://localhost:6001/users/budmore \
    >   --header 'authorization: valid-token'

    {"username":"budmore","role":"meetup"}

Greeting user:

    $ curl --request POST \
    >   --url http://localhost:6001/users/budmore/greet \
    >   --header 'authorization: valid-token' \
    >   --header 'content-type: application/json' \
    >   --data '{"greeting":"hola"}'
    
    {"greeting":"Hola, budmore! You are a great meetup!"}

It works! It's that simple! The full source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/grpc).
