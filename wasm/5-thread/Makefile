GOROOT=$(shell go env GOROOT)

build:
	mkdir -p dist
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" dist/wasm_exec.js
	cp "$(GOROOT)/misc/wasm/wasm_exec.html" dist/index.html
	GOARCH=wasm GOOS=js go build -o dist/test.wasm wasm.go

clean:
	rm -rf dist/*

run:
	go run main.go
