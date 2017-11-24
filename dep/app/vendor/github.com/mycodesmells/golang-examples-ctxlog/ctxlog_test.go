package ctxlog

import (
	"context"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	ctx := context.WithValue(context.Background(), "label", "depdep")
	Println(ctx, time.Now(), "Hello world!")
}
