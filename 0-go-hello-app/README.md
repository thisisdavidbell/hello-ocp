# 0 - Creating a Go hello world application in Docker

This repo guess you through the process of creating a simple Hello World httpserver go app. You will then produce a docker container running the application. Finally, you will use some of the build in Red Hat Openshift Container Platform tools to quickly and easily run the application in docker. 

In later chapters, this simple application will be the one managed by the operator you create.

## 0. Prereqs

- an OCP 4.4 or above cluster.
  - The tutorial has been tested on a range of OCP providers, and at OCP 4.4 and 4.5.
- docker
- You must overwrite key values in the examples
  - This tutorial requires a number of values which are specific to your setup. Therefore, in all sections of this tutorial, you should search and replace the following values appropriately:
    - DOCKERHOSTNAME - the fully qualified hostname of your docker repo
    - DOCKERNAMESPACE - the namespace within your docker repo where all images will be pushed.
    - OCPHOSTNAME - the fully qualified hostname of your ocp system. (e.fg. everything after `api.` or `apps.` )

## 1. Create Go app
Create a simple go application, running a httpserver, returning a string.
You can use the simple example here: [hello-ocp.go](hello-ocp.go)

## 2. Test it
In one terminal, run:
`go run hello.go`

In a second terminal run:
`curl localhost:8080/hello`

You can now stop the app in the first terminal.

## 3. Build it
Run:
`go build hello.go`

## 4. Run it
Run and test it:
`./hello`
`curl localhost:8080/hello`
Stop the app

## 5. Create a dockerfile
Red Hat Openshift Container Platform utilises docker containers. You should therefore write a dockerfile to wrap the go application.
An example dockerfile is provided: [Dockerfile](dockerfile)
Test the dockerfile:
- Build the image
`docker build -t hello:v0.0.1 .`

- Confirm the image was created
`docker images | grep hello`

- Run the image locally:
`docker run -p 8080:8080 -d hello:v0.0.1`

- Confirm the container is running:
`docker ps`

- Test:
`curl localhost:8080/hello`

- Stop the container:
`docker stop CONTAINER-ID`
 where `CONTAINER-ID` is the value shown when running `docker ps`

## 6. Create a Makefile
It is common to use Makefile's or similar technologies to group together the common commands used. 
A simple helper [Makefile](Makefile) is provided covering a number of commands mentioned above, and introduced later.

## 7. Create new OCP Application in CRC
While not required for this tutorial, it is interesting to use some of the OCP tooling to quickly and easily deploy this application in OCP.

This step will show how to do both the following
 - a. creating an app in OCP UI from an existing Dockerfile
 - b. creating app from CLI from source code

### 8.a Creating an app in OCP UI from an existing Dockerfile
 - Log into your OCP web console through your browser of choice.
 - Select the 'Developer' perspective
 - Create a new project, for example `hello-dockerfile`
 - Click `+Add` In left hand menu
 - Select `From Dockerfile`
 - Ensure all values are correct. Specifically:
    - git repo url, e.g. `https://github.com/thisisdavidbell/hello-ocp`
    - Container port - enter the port specified by `EXPOSE` in the Dockerfile, e.g. 8080
    - Resources - select DeploymentConfig for more Openshift specific functionality
    - Create a Route - leave ticked
 - Click on `Routing`
   - enter a hostname, including the full hostname of your OCP system, e.g. hello-dockerfile.apps.OCPHOSTNAME
   - Path: `/hello`
   - Target port - enter the same port as above, e.g. 8080 (A service will be created which exposes this port
 - Click `Create`
 - In the OCP web console, view the build, the deploymentConfig, service and route
 - Test the application
   - Select the 'Administrator' perspective
   - Networking
   - Routes
   - Click the link under 'Location' for the appropriate route.
     - Note this is just a http url. You can also use `curl URLFROMLOCATIONFIELD`
     - Note: if this fails initially, the pod running your application may not be up yet. Try again in a minute.

### 8.b Create app from CLI from source code
// _TODO_ - test later once pushed to public repo!!!

 - ensure you are connected to an OCP cluster
   - `oc login`
 - in root dir of hello-ocp repo, run:
    - `oc new-project hello-sourcecode`
    - `oc new-app .`

Amazingly, that is all you need to do.
OCP will now go off and spot this is go code, build a go image, push that into the internal image registry in OCP, create image streams, etc and deploy the image as a DeploymentConfig. It didn't however create a route.

- Create a Route
// _TODO_ add in CLI method to create this.
 - perform route UI step above, only with the host: `hello-sourcecode.app.OCPHOSTNAME/hello`
 
 - Test the application:
   `curl  hello-sourcecode.app.OCPHOSTNAME/hello`
   Note: you didn't specify a port, so the http default port of 80 is used. A route services the port up on 80 by default

See the Appendix at the bottom of this readme for rebuilding the image.

## 9. Push image to docker
You will need the image in a docker registry accessible from your OCP system for future sections.
`docker tag hello:v0.0.1 DOCKERHOSTNAME/DOCKERNAMESPACE/hello:v0.0.1`
`docker push DOCKERHOSTNAME/DOCKERNAMESPACE/hello:v0.0.1`