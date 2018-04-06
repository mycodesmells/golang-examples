FROM golang:1.10-stretch

COPY . ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/blog-admin

WORKDIR ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/blog-admin

RUN make build/docker

# end of first stage, beginning of the second one
FROM alpine:3.7

COPY --from=0 /blog-admin /blog-admin
CMD ["/blog-admin"]