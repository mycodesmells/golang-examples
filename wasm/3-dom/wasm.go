// +build js,wasm

package main

import (
	"syscall/js"
)

func main() {
	registerCallbacks()
	c := make(chan struct{})
	<-c
}

func registerCallbacks() {
	js.Global().Set("write", js.NewCallback(write))
	js.Global().Set("createDiv", js.NewCallback(createDiv))
}

func write(args []js.Value) {
	js.Global().Get("document").Call("write", args[0])
}

func createDiv(args []js.Value) {
	document := js.Global().Get("document")

	divText := document.Call("createTextNode", args[0])
	div := document.Call("createElement", "div")
	div.Set("style", "background: tomato;")
	div.Call("appendChild", divText)

	document.Get("body").Call("appendChild", div)
}
