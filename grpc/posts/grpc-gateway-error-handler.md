# gRPC gateway Error Handler

Before going any further make sure you've read a previous post, [about grpc-gateway](http://mycodesmells.com/post/http-server-from-grpc-with-gateway).

Exposing gRPC server in a form of HTTP API is pretty easy using grpc-gateway, but when it comes to returning errors, the default behavior is something that might be bothering us. Let's see how can we improve that.

### Default error handler

When we set up our gRPC-based HTTP server we defined some basic level of authentication, as we required `Authorization` header with a value of `valid-token`. When we don't provide it, we get a correct, 401 status response, but when we take a look at the response body, we get just a bit more that we'd like:

    {
        "error": "invalid token",
        "code": 16
    }

Although it's pretty cool, that gRPC's error code 16 (_Unauthenticated_) is correctly translated to HTTP's 401, the original one is still returned. This results in our API revealing that we are running gRPC under the sheets. We would definitely like to avoid that, and fortunately, we can.

When we dig a bit deeper in the codebase of [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) we can find that whenever an error is returned from the gRPC handler, it goes through a specific error handler stored in variable `HTTPError` (in `github.com/grpc-ecosystem/grpc-gateway/runtime` package). Its default implementation gives us an idea of what our custom handler should look like:

    // github.com/grpc-ecosystem/grpc-gateway/runtime/errors.go
    ...
    
    ...
    func DefaultHTTPError(ctx context.Context, mux *ServeMux, marshaler Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
        ...
    }

### Custom handler

Our custom handler needs to do three things: return a correct HTTP status, a correct `Content-type` header and a body that contains only the message port of the error.

We start by `Content-type` which is the simplest as an input parameter, `runtime.Marshaler`, provides an utility function called... `ContentType()`:

    func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
        w.Header().Set("Content-type", marshaler.ContentType())
        ...
    }

Then we take care of the status code. This is slightly more complicated, but after all, we just need to copy that from the default handler. First, we need to get a gRPC error code from `err` (`grpc.Code(err)`), then translate it to HTTP status code (`runtime.HTTPStatusFromCode(..)`) and pass it to the response writer.

    ...
    w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
    ...

Last, but certainly not least, we need to build a response body. Let's create a simple struct for that:

    type errorBody struct {
        Err string `json:"error,omitempty"`
    }

Now, all we need to to is insert our error message into that struct and encode it to JSON. You might want to use `err.Error()` function for that:

    ...
    json.NewEncoder(w).Encode(errorBody{
        Err: err.Error(),
    })
    ...

In order to see our handler in action, we need to overwrite `HTTPError` in `runtime` package and restart the server. Unfortunately, when we call our API now, we realize that we've taken a step back:

    {
    "error": "rpc error: code = Unauthenticated desc = invalid token"
    }

That is because we should not use `err.Error()` as it includes some extra information revealing that it is, in fact, a gRPC error. Instead, we need to use `grpc.ErrorDesc(err)` which gets rid of all that, leaving only a true message which we want to return to the API user:

    {
    "error": "invalid token"
    }

As you see, making our error responses cleaner did not require us too much effort. The handler, with some fallback logic in case of JSON marshaling errors, looks like this:

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

The full source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/grpc).
