// +build js,wasm

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			fmt.Printf("#%d\n", n)
			time.Sleep(time.Second)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
