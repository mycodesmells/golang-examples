package log

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	Println(time.Now(), "Hello world!")
}
