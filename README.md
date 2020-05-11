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

## 6. Create a new publicly accessible git repo

## 7. Push these files to the repo

## 8 Create new OCP Application in CRC

For purposes of this demo, choose between:
 - a. creating app in UI from Dockerfile
 - b. creating app from CLI from source code

### 8.a Create app from UI from Dockerfile

 - launch OCP console: `crc console`
 - Developer perspective
 - Click `+Add` In left hand menu
 - Select `From Dockerfile`
 - Ensure all values are correct. Specifically:
    - git repo url, e.g. `https://github.com/thisisdavidbell/hello-ocp`
    - Container port - enter the port specified by `EXPOSE` in the Dockerfile
    - Resources - select DeploymentConfig for more Openshift specific functionality
    - Create a Route - leave ticked
 - Click on `Routing`
   - enter a hostname, including the default dns of `apps-crc.testing`, e.g. hello-ocp.apps-crc.testing
   - Path: `/hello-ocp`
   - Target port - 8080 (A service will be created which exposes this port (_TODO_ - I already entered the container port, so is this the port the service exposes, that the Route then points at? ))
 - Click `Create`

### 8.b Create app from CLI from source code

 - ensure you are connect to an OCP
 - in root dir of hello-ocp repo, run:
    - `oc new-project project1`
    - `oc new-app .`

Amazingly, that is all you need to do.
OCP will now go off and spot this is go code, build a go image, push that into the internal image repo in OCP, create image streams, etc and deploy the image as a DeploymentConfig.

What it doesn't do is create a route, so we should do that.
 - perform route UI step above

 _TODO_ add in CLI method to create this.

## 9. Test the app

Once the build process has finished, you should now be able to run the app:
 - Run: `curl  hello-ocp.apps-crc.testing/hello-ocp`


## 10. View the build etc

_TODO_

## 11. View the deploymentConfig, Service, Route in the console

_TODO_

## 12. Make a change to the app, rebuild the image

- Make a change to the go app, e.g. change the Hello message.
- Push the change to the public repo, or just change locallt.
- Remember the git commit message

- Select Administrator Perspective
- Select Builds->BuildConfigs
- Select Rebuild
- View the build under Builds
- Note the newly running build shows the new commit message.

_TODO_ check this works for CLI local source approach too

## 13. Test the app

Once the build process has finished, you should now be able to run the app again:
 - Run: `curl  hello-ocp.apps-crc.testing/hello-ocp`
 - Note the new message is now returned.

_TODO_: Can we update the triggers to automatically spot this (if it wouldn't already given time)?
