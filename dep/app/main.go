package main

import (
	"context"
	"time"

	ctxlog "github.com/mycodesmells/golang-examples-ctxlog"
	log "github.com/mycodesmells/golang-examples-log"
)

func main() {
	ctx := context.WithValue(context.Background(), "label", "dep-app")

	log.Println(time.Now(), "Hello world from log")
	ctxlog.Println(ctx, time.Now(), "Hello world from ctxlog")
}
