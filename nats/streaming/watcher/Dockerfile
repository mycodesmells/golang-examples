FROM golang:1.10-stretch

COPY . ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/streaming/watcher
WORKDIR ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/streaming/watcher

RUN make build/docker

# end of first stage, beginning of the second one
FROM alpine:3.7

COPY --from=0 /watcher /watcher
CMD ["/watcher"]