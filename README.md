# hello-ocp - Creating a hello world operator in Red Hat Openshift Container Platform (OCP)

This repo provides a step by step tutorial using the operator-sdk to create an Red Hat Openshift Container Platform (OCP) operator which manages a simple hello world application.

## Tutorials

Each section below can either just be read and the example code reviewed, or followed step by step as you perform every step yourself. Any section can be performed in isolation, but for the best experience start at the beginning and work through each section in turn.

Important notes:
- This tutorial requires a number of values which are specific to your setup. Therefore, in all sections of this tutorial, you should search and replace the following values appropriately:
- DOCKERHOSTNAME - the fully qualified hostname of your docker repo
- DOCKERNAMESPACE - the namespace within your docker repo where all images will be pushed.
- OCPHOSTNAME - the fully qualified hostname of your ocp system. (e.fg. everything after `api.` or `apps.` )

- [0 - Creating a Go hello world application in Docker](0-go-hello-app)
- [1 - Creating the operator](1-create-operator)

