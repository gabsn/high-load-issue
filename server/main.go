package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	}))
}
