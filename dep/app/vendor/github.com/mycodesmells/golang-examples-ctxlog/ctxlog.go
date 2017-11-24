package ctxlog

import (
	"context"
	"fmt"
	"time"

	log "github.com/mycodesmells/golang-examples-log"
)

func Println(ctx context.Context, t time.Time, v ...interface{}) {
	label, ok := ctx.Value("label").(string)
	if ok && label != "" {
		v = append(v, 0)
		copy(v[1:], v[:])
		v[0] = fmt.Sprintf("[label=%s]", label)
	}
	log.Println(t, v...)
}
