package main

import (
	"fmt"
	"net/http"
)

var PORT = "8080"

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello yet again OCP.\n")
}

func main() {
	fmt.Printf("Running hello server on %s\n", PORT)
	http.HandleFunc("/hello-ocp", hello)
	http.ListenAndServe(":"+PORT, nil)
}
