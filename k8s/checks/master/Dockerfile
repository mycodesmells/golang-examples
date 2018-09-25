FROM golang:1.11-alpine3.8

COPY . /go/src/github.com/mycodesmells/golang-examples/k8s/checks/master
RUN go install github.com/mycodesmells/golang-examples/k8s/checks/master

FROM alpine:3.8
RUN apk add --no-cache ca-certificates

COPY --from=0 /go/bin/master /master
CMD ["/master"]
