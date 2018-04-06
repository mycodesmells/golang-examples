test:
	go test -race ./...

go-bindata:
	go-bindata -o templates.go *.tpl

build/docker: test go-bindata
	CGO_ENABLED=0 go build -a -installsuffix cgo -o /blog-generator .

