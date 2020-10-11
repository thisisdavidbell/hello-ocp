#Make file to create and deploy hell-ocp

.PHONY: operator-build-and-push
operator-build-and-push:
	operator-sdk build somerandomhostnamem/drb/hello-ocp-operator:v0.0.1
	docker push somerandomhostnamem/drb/hello-ocp-operator:v0.0.1

.PHONY: operator-deploy
operator-deploy:
	echo TODO