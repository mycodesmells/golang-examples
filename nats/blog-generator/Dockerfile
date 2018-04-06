FROM golang:1.10-stretch

COPY . ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/blog-generator
WORKDIR ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/blog-generator

RUN go get -u github.com/kevinburke/go-bindata/...
RUN make build/docker

# end of first stage, beginning of the second one
FROM alpine:3.7

COPY --from=0 /blog-generator /blog-generator
CMD ["/blog-generator"]