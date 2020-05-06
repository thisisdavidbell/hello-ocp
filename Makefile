#Make file to create and deploy hell-ocp

.PHONY: all
all: go-compile docker-build

.PHONY: go-compile
go-compile:
	go build hello-ocp.go

.PHONY: go-run
go-run:
	./hello-ocp

.PHONY: docker-build
docker-build:
	docker build -t hello-ocp .

.PHONY: docker-run
docker-run:
	docker run -p 8081:8080 -d hello-ocp:latest
