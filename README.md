README

This repo tracks the process of creating a simple Hello World httpserver go app, and the process required to deploy it into Openshift Containers Platform. It was created using CodeReady Containers (crc).

## 1. Create Go app
Simple go app here: [hello-ocp.go](hello-ocp.go)

## 2. Test it
Run:
`go run hello-ocp.go`

## 3. Build it
Run:
`go build hello-ocp.go`

## 4. Run it
Run:
`./hello-ocp`

## 5. Create a Makefile
A simple helper Makefile: [Makefile](Makefile)
For example, this allows:
- `make go-compile`
- `make go-run`
- `make docker-build`
- `make docker-run`
