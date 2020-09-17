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

.PHONY: operator-build
operator-build:
	operator-sdk build somedockerrepohostname/drb/hello-ocp-operator:v0.0.1
	docker push somedockerrepohostname/drb/hello-ocp-operator:v0.0.1

.PHONY: operator-deploy
operator-deploy:
	echo TODO