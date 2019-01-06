// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	registerCallbacks()
	c := make(chan struct{})
	<-c
}

func registerCallbacks() {
	js.Global().Set("hello", js.NewCallback(sayHello))
}

func sayHello(i []js.Value) {
	shout := fmt.Sprintf("Hello, %s! Welcome to the GoWroc talk!", i[0])
	js.Global().Get("console").Call("log", shout)
}
