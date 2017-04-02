# Authentication in gRPC

Setting up communication between the server and its clients via gRPC is really simple, but what about the security? It's not that hard, either. Let's take a look.

### Insecure is awful

First of all, there is something that has to be said about the security in general. If you don't keep your communication secure, you are doing it wrong. On the [GoLab 2017](https://golab.io/) conference I saw a short, but powerful lightning talk delivered by Eleanor McHugh [feyeleanor](https://twitter.com/feyeleanor) during which she said that anything that you would not like to see printed and handed out in every public place around the world is what you should keep secure. Powerful, right? What is your argument for not using any encryption?

### Certificates and keys

Just as with HTTP servers, a gRPC server needs a certificate and key pair in order to have its communication encrypted. Since we are developing a local example, we don't need to have it delivered by any official CA (_Certificate Authority_) - we can generate them on our own. In order to do this we use `openssl` command:

    openssl req -x509 -newkey rsa:4096 -keyout cmd/server/server-key.pem -out cmd/server/server-cert.pem -days 365 -nodes -subj '/CN=localhost'

Besides the paths to keep the output files we need to provide the expiration time (365 days), algorithm details (RSA, 4096-bit length), format standard (X.509) and our server's common name (CN) which is just a `localhost` in our case.

### Usage in gRPC

In order to run our server in a more secure way, we need to read our certificate and key from the disk and create so caled TLS credentials:

    // cmd/server/main.go
    ...
    creds, err := credentials.NewServerTLSFromFile("cmd/server/server-cert.pem", "cmd/server/server-key.pem")
    if err != nil {
        log.Fatalf("Failed to setup tls: %v", err)
    }
    ...
    
Then we need to pass those `creds` to our server:

    // cmd/server/main.go
    ...
    server := grpc.NewServer(
        grpc.Creds(creds),
    )
    pb.RegisterSimpleServerServer(server, NewServer())
    
    server.Serve(lis)
    ...

If we now try to connect to the server using our old client, we are rejected pretty quickly:

    $ go run cmd/client/main.go 
    2017/04/02 23:27:42 transport: http2Client.notifyError got notified that the client transport was broken unexpected EOF.
    2017/04/02 23:27:42 Failed to create user: rpc error: code = Internal desc = transport is closing
    exit status 1

While this looks a bit mysterious, the other side of the equation reveals what is the problem:

    $ go run cmd/server/main.go 
    2017/04/02 23:27:38 Listening on :6000
    2017/04/02 23:27:42 grpc: Server.Serve failed to complete security handshake from "[::1]:61133": tls: first record does not look like a TLS handshake

Fortunately, implementing the _TLS handshake_ is also easy on the client's side. We only need to import the appropriate certificate file:

    // cmd/client/main.go
    creds, err := credentials.NewClientTLSFromFile("cmd/server/server-cert.pem", "")
    if err != nil {
        log.Fatalf("cert load error: %s", err)
    }
    ...
    conn, err := grpc.Dial("localhost:6000", grpc.WithTransportCredentials(creds))
    ...
    
And now the output looks much better:

    $ go run cmd/client/main.go 
    2017/04/02 23:30:56 Created user!
    2017/04/02 23:30:56 User exists: username:"slomek" role:"joker" 
    2017/04/02 23:30:56 Greeting: Howdy, slomek! You are a great joker!
    

### Additional security

While using certificates are great, it's hard to authenticate which exact user (client) is making requests to your gRPC server, right? After all, the certificate is common for all clients. In order to recognize who is using our service, we need to make use of the context and store some data there. In order to make it as clean as possible, we can use something called `UnaryInterceptor` which allows us to block some gRPC calls if something is wrong with them. Such interceptor is a simple function:

    func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        ...
    }

So the simplest approach would be to put some value in the context on one side and read it in the interceptor, right? Let's do it then:

    // cmd/client/main.go
    ...
    ctx := context.WithValue(context.Background(), "token", "valid-token")
    
    _, err = client.CreateUser(ctx, &pb.CreateUserRequest{User: &pb.User{Username: "slomek", Role: "joker"}})
    ...

And read it on the other side:

    // cmd/server/main.go
    ...
    func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        t := ctx.Value("token")
        log.Printf("Token: %v", t)
        if t == nil {
            return nil, grpc.Errorf(codes.Unauthenticated, "incorrect access token")
        }
        token := t.(string)
        if token != "valid-token" {
            return nil, grpc.Errorf(codes.Unauthenticated, "incorrect access token")
        }
        
        return handler(ctx, req)
    }
    ...
    server := grpc.NewServer(
        grpc.Creds(creds),
        grpc.UnaryInterceptor(AuthInterceptor),
    )
    ...
    
Let's see if this works:

    $ go run cmd/client/main.go 
    2017/04/02 23:46:45 Failed to create user: rpc error: code = Unauthenticated desc = incorrect access token
    exit status 1

What? Let's check what does server say:

    $ go run cmd/server/main.go 
    2017/04/02 23:46:42 Listening on :6000
    2017/04/02 23:46:45 Token: <nil>

The thing is, that the `context` does not pass the values from one end of gRPC communication to another. Is there anything we can do in this case? Of course, it is! There is a special package in `"google.golang.org/grpc/metadata"` that provides us with this feature. Using this `metadata` we need to create our context in a special way:

    // cmd/client/main.go
    ...
    md := metadata.Pairs("token", "valid-token")
    ctx := metadata.NewContext(context.Background(), md)
    
    _, err = client.CreateUser(ctx, &pb.CreateUserRequest{User: &pb.User{Username: "slomek", Role: "joker"}})
    ...
    
Reading it is also slightly different. We need to check if a metadata exists, then check it our key is present. Each metadata value is a slice, so we need to make sure that our slice is not empty (and contains one and only one element), and contains expected value:

    ...
    func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        meta, ok := metadata.FromContext(ctx)
        if !ok {
            return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
        }
        if len(meta["token"]) != 1 {
            return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
        }
        if meta["token"][0] != "valid-token" {
            return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
        }
        
        return handler(ctx, req)
    }
    ...
    
Now this should finally work:

    $ go run cmd/client/main.go
    2017/04/02 23:52:52 Created user!
    2017/04/02 23:52:52 User exists: username:"slomek" role:"joker" 
    2017/04/02 23:52:52 Greeting: Howdy, slomek! You are a great joker!
    
    $ go run cmd/server/main.go 
    2017/04/02 23:52:50 Listening on :6000
    2017/04/02 23:52:52 Creating user...
    2017/04/02 23:52:52 User created!
    2017/04/02 23:52:52 Getting user!
    2017/04/02 23:52:52 User found!
    2017/04/02 23:52:52 Greeting user...
    2017/04/02 23:52:52 Getting user!
    2017/04/02 23:52:52 User found!
    
It does! The full source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/grpc).
