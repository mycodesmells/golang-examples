test:
	go test -race ./...

build/docker: test
	CGO_ENABLED=0 go build -a -installsuffix cgo -o /blog-admin .
