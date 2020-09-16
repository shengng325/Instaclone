package main

import "net/http"

func main() {
	println("Hello, World!")
	http.ListenAndServe(":3000", nil)
}
