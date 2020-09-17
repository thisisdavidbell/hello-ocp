package main

import (
	"fmt"
	"net/http"
	"os"
)

// PORT - server port
var PORT = "8080"

// PATH - server path
var PATH = "/hello-ocp"

func hello(w http.ResponseWriter, req *http.Request) {
	helloName := "helloNameEnvVarNotSet"
	envName := os.Getenv("HELLONAME")
	if envName != "" {
		helloName = envName
	}
	fmt.Fprintf(w, "hello %s. Im v0.0.1.\n", helloName)
}

func main() {
	fmt.Printf("Running hello server on %s:%s\n", PATH, PORT)
	http.HandleFunc(PATH, hello)
	http.ListenAndServe(":"+PORT, nil)
}
