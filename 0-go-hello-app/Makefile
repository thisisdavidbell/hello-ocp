#Make file to create and deploy hell-ocp

.PHONY: all
all: go-build docker-build

.PHONY: go-build
go-compile:
	go build hello.go

.PHONY: go-run
go-run:
	./hello

.PHONY: docker-build
docker-build:
	docker build -t hello .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -d hello:v0.0.1

