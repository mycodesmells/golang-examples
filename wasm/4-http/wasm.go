// +build js,wasm

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"
)

func main() {
	urlInput := js.Global().Get("document").Call("getElementById", "urlInput")
	url := urlInput.Get("value").String()

	ping(url)
	// WORKS: https://httpbin.org/anything
}

func ping(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to ping '%s': %v", url, err)
		return
	}
	defer resp.Body.Close()

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response from '%s': %v", url, err)
		return
	}
	fmt.Printf("Response from '%s':", url)
	fmt.Println(string(bb))
}
