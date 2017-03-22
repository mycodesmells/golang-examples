gen_proto:
	protoc --go_out=plugins=grpc:${GOPATH}/src proto/service/service.proto
	protoc --go_out=${GOPATH}/src proto/message/message.proto

run:
	go run main.go

run/server:
	go run cmd/server/main.go

run/client:
	go run cmd/client/main.go
