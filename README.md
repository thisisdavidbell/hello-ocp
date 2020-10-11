# hello-ocp - Creating a hello world operator in Red Hat Openshift Container Platform (OCP)

This repo provides a step by step tutorial using the operator-sdk to create an Red Hat Openshift Container Platform (OCP) operator which manages a simple hello world application.

# Tutorials

Each section below can either just be read and the example code reviewed, or followed step by step as you perform every step yourself. Any section can be performed in isolation, but for the best experience start at the beginning and work through each section in turn.

- [0 - Creating a Go hello world application in Docker](0-go-hello-app)
- [1 - Creating the operator](1-create-operator) _WIP_

Done but doc not yet updated:
- [2 - Add some reconcile logic](2-reconcile) // hello name in cr
- [3 - Create catalog, bundle images and catalog source](3-catalog-bundle-catsrc)
- [4 - Operator managed Service and Route](4-service-route)
- Add validation via annotations

ToDo:
- Update to create a deployment, and use size and other fields
- Create a new version of operator, and a channel and have a currently deployed operator auto updated
- Update operator to add route (service already done)
- Implement spec.version to allow user to specify which version of hello they want (maps to image tag)
- Add a validating webhook using operator-sdk 1.0
- Update Status section of of cr with useful values
  - url of helloocp endpoint would be great start. Perhaps previous name used another status
- Add a logo and description to operator hub