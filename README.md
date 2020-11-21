# hello-ocp - Creating a hello world operator in Red Hat Openshift Container Platform (OCP)

This repo provides a step by step tutorial using the operator-sdk (version 1.0) to create a Kubernetes operator which manages a simple hello world application. We will run this operator on Red Hat Openshift Container Platform (OCP).

# Tutorial

The tutorial is written so that you can start from scratch, and work through every chapter in turn. You will end up with a feature rich operator managing the lifecycle of a simple hello world go application. Alternatively, you can just read each section and view the code that would be produced at the end of the section by looking in each corresponding sub-directory in the [checkpoints](checkpoints) directory. You can also use the appropriate checkpoint to start the tutorial at any chapter.

- [1 - Creating a Go hello world application in Docker](README-1-hello-go-app.md)
- [2 - Creating the operator](1-create-operator) _WIP_
- [3 - Deploy operator ]() - look at tutorial/my notes to just build operator and push to ocp
- [4 - Deploy operator with olm]() - look at tutorial/my notes to build csv and push to OCP
- [5 - Deploy operator to operator hub]() - use [previous](3-catalog-bundle-catsrc))
- [6 - make updates to hello application in reconcile loop] 
  - create service in reconcile loop - see previous code as already done - use [previous](4-service-route)
  - create route in reconcile loop
  - [hello name and version in cr - see previous code as already done name]() - use [previous](2-reconcile)
  - add web hook using operator-sdk 1.0


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
- Update Status section of cr with useful values
  - url of helloocp endpoint would be great start. Perhaps previous name used another status
- Add a logo and description to operator hub