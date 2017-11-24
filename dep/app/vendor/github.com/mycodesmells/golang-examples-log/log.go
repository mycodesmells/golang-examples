package log

import (
	"fmt"
	"time"
)

func Println(t time.Time, v ...interface{}) {
	v = append(v, 0)
	copy(v[1:], v[:])
	v[0] = t.Format(time.RFC3339)
	fmt.Println(v...)
}
