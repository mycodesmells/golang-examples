.PHONY: all
all:
	@echo "usage: 'make build'"

.PHONY: build
build: assets
	go build -o statik-example .

.PHONY: run
run: assets
	go run .

.PHONY: assets
assets:
	statik -include='*.txt' -src=.
