package main

import (
	"fmt"
	"net/http"
)

var PORT = "8080"

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello OCP. Im v0.0.1.\n")
}

func main() {
	fmt.Printf("Running hello server on %s\n", PORT)
	http.HandleFunc("/hello-ocp", hello)
	http.ListenAndServe(":"+PORT, nil)
}
