FROM golang:1.10-stretch

COPY . ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/streaming/neatflyx

WORKDIR ${GOPATH}/src/github.com/mycodesmells/golang-examples/nats/streaming/neatflyx

RUN make build/docker

# end of first stage, beginning of the second one
FROM alpine:3.7

COPY --from=0 /neatflyx /neatflyx
CMD ["/neatflyx"]