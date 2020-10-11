# 1 - Creating the operator

## 0. Prereqs

- a dockerised go application, with image pushed to accessible docker egistry. You can complete section 0 to achieve this.
- an OCP 4.4 or above cluster.
  - The tutorial has been tested on a range of OCP providers, and at OCP 4.4 and 4.5.
- docker running locally
- operator-sdk installed - https://sdk.operatorframework.io/docs/installation/install-operator-sdk/. Note, this tutorial was originally writted using operator-sdk 0.17.1. However, operator-sdk version 1.0 will be needed to create the web hook in later sections.
- You must overwrite key values in the examples
  - This tutorial requires a number of values which are specific to your setup. Therefore, in all sections of this tutorial, you should search and replace the following values appropriately:
    - DOCKERHOSTNAME - the fully qualified hostname of your docker repo
    - DOCKERNAMESPACE - the namespace within your docker repo where all images will be pushed.
    - OCPHOSTNAME - the fully qualified hostname of your ocp system. (e.fg. everything after `api.` or `apps.` )

## Useful resources:
It is highly recommended you read through at least one example of creating an operator. The etcd examples in either ofthe following 2 sources are a great start to understanding operators. This tutorial assumes you have aleady build up some level of knowledge of OCP and operators, and so have some understanding of terms such as pod, service, route, operator, operator-sdk, reconcile loop, bundle image, catalog image, catalog source.

- Kubernetes Operators O'Reilly book by Jason Dobies & Joshua Wood - https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/ 
  - this is a fantastic introduction to operators, and will provide all the background knowledge you need. Note however that it uses the deprecated package manifest approach to packaging rather than the newer recommended bundle approach. We will use the bundle aproach in this tutorial.
- https://operatorframework.io/ - A great source of help and tutorials
  - https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/ this is the latest version of the more indepth go operator tutorial, creating a basic operator for the existing etcd container.

## Useful notes:
- OLM is installed by default on OCP 4.4 and higher, so you do not need to manually install this.

## 1. Create an operator to manage the hello application

### Create new directory
The files present in the 1-create-operator directory are the result of completing all the steps of this section.
To perform all the steps yourself, create a new directory inside your $GOPATH/src directory structure and cd into it.
Ensure go module support is on by setting:
`export GO111MODULE=on`

### Initialise the project
Create your initial operator project
`operator-sdk init`
note: this used the command `new` in earlier versions of operator-sdk`

_WIP - HERE_

### Create your resource definition and controller
You need to define a new resource (CRD deascribing the custom resource yaml), and a controller (the operator code which manages the hello application)

`WIP - read up on new command params on cli - operator-sdk create api --version=v1alpha1 --kind=hello`
note: this used the command `add` in earlier versions of operator-sdk`



Following steps of linked tutorial, with following changes:
 - _TODO_ insert all steps here to allow following along here only
 - command to create crd: `operator-sdk add api --api-version=helloocp.example.com/v1alpha1 --kind=Helloocp`

 - adding an extra spec and status item to see how its handled.
   -  Note, these I believe will be completely ignored as we dont write any code to use them (size makes no sense for a pod anyway).

 - updated crd helloocp_types.go with:
 ```
 type HelloocpSpec struct {
 	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
 	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
 	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

  // +kubebuilder:validation:Maximum=3
	// Size is the size of the memcached deployment
	Size int32 `json:"size"`

	// +kubebuilder:validation:Enum=Option1;Option2
	// An extra spec field - a string to see what happens
	SomeString string `json:"someString"`
 }

 // HelloocpStatus defines the observed state of Helloocp
 type HelloocpStatus struct {
 	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
 	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
 	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

 	// Nodes are the names of the memcached pods
 	Nodes []string `json:"nodes"`

 	// An extra status field - a string to see what happens
 	VersionString string `json:"versionString"`
 }
 ```
 - NOTE: I added validation too.
 - ran: `operator-sdk generate crds` to update crd with new validation
 - added `  someString: imSomeString` to end or cr file (`hello-ocp-operator/deploy/crds/helloocp.example.com_v1alpha1_helloocp_cr.yaml`). It didnt get added, so assuming I need to manually (?)

 - createcontroller command: `operator-sdk add controller --api-version=helloocp.example.com/v1alpha1 --kind=Helloocp`

 - update PodSpec:
 ```
 Spec: corev1.PodSpec{
   Containers: []corev1.Container{
     {
       Name:    "busybox",
       Image:   "busybox",
       Command: []string{"sleep", "3600"},
     },
   },
```
replacing busybox with with image, name and specify command (we will push image to internal registry shortly.) so:
```
Spec: corev1.PodSpec{
  Containers: []corev1.Container{
    {
      Name:    "hello-ocp",
      Image:   "image-registry.openshift-image-registry.svc:5000/project1/hello-ocp:v0.0.1",
      Command: []string{"./hello-ocp"},
    },
  },
},
  ```



- I updated the msg printed by hello-ocp.go to include the string `v0.0.1`
- `docker build -t hello-ocp:v0.0.1 .`
- `docker tag hello-ocp:v0.0.1 default-route-openshift-image-registry.apps-crc.testing/project1/hello-ocp:v0.0.1`
- if not already done:
  - expose registry on default route:
    - `oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge`
  - grant permission for kubeadmin to use registry:
    - `oc policy add-role-to-user registry-viewer kubeadmin`
  - register the internal registry as an insecure registry in docker (follow this or equivalent):
    - Docker->Preferences->daemon-> add default-route-openshift-image-registry.apps-crc.testing as insecure registry
- `docker login -u kubeadmin -t <oc whoami -t> default-route-openshift-image-registry.apps-crc.testing/project1/hello-ocp`
- `docker push default-route-openshift-image-registry.apps-crc.testing/project1/hello-ocp:v0.0.1`

- register crd: `oc create -f deploy/crds/helloocp.example.com_helloocps_crd.yaml`
  - note, if adding validation to types file `helloocp_types.go` and regenratoing crd to include the validation, you only need to redeploy the crd file
- build operator: `operator-sdk build somerandomhostnamem/drb/hello-ocp-operator:v0.0.1`
- push image: `docker push somerandomhostnamem/drb/hello-ocp-operator:v0.0.1`
- update `image:`  in deploy/operator.yaml to `image: image-registry.openshift-image-registry.svc:5000/project1/hello-ocp-operator:v0.0.1`
- deploy as instructed to create operator without OLM:
```
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml
```
- create a 'hello-ocp': `oc apply -f deploy/crds/helloocp.example.com_v1alpha1_helloocp_cr.yaml`
- confirm pod up and running
  - _TODO_ I noted the pod log was empty - no `server running on 8080` message - why?
- set up a service manually, including values:
```
port: 8080
targetPort: 8080
selector:
  app: example-helloocp-3
```
- set up a route manually
  - Name: `example-helloocp-1-route`
  - Hostname: `example-helloocp-1.apps-crc.testing`
  - Path: `/hello-ocp`

  - test: `curl http://hello-ocp.apps-crc.testing/hello-ocp`
  - SUCCESS


### Deploying Helloocp operator with OLM

Following: https://github.com/operator-framework/getting-started#manage-the-operator-using-the-operator-lifecycle-manager

- ensured deleted items from operator only tutorial above.
- created csv using: `operator-sdk generate csv --csv-version 0.0.1 --update-crds`
- note new csv and crd files appeared in `deploy/olm-catalog/hello-ocp-operator/manifests`
- need to jump to remaining steps in: https://docs.openshift.com/container-platform/4.4/operators/operator_sdk/osdk-getting-started.html#managing-memcached-operator-using-olm_osdk-getting-started in order to deploy OLM managed operator:
  - update `namespace:`` to `project1`
    - _TODO_: understand why this is required.
  - create an operator group contaiing namespace to deploy into:
```
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: memcached-operator-group
  namespace: default
spec:
  targetNamespaces:
  - drb-hello
```
NOTE: update namespace from `default` to `project1`
 - deploy the ClusterServiceVersion: `oc apply -f deploy/olm-catalog/hello-ocp-operator/manifests/hello-ocp-operator.clusterserviceversion.yaml`
 - `oc get deployments`
 - deploy the crd (note the 2 crds from operator and OLM are identical): `oc apply -f deploy/olm-catalog/hello-ocp-operator/manifests/helloocp.example.com_helloocps_crd.yaml`
- Copy cr yaml: `cp deploy/crds/helloocp.example.com_v1alpha1_helloocp_cr.yaml helloocp-cr.yaml`
- update to fail the validation:
```
apiVersion: helloocp.example.com/v1alpha1
kind: Helloocp
metadata:
  name: example-helloocp-3
spec:
  # Add fields here
  size: 4
  someString: imSomeString
```
- Deploy: `oc apply -f helloocp-cr.yaml`
- Note the failure due to both validation rules failing
- Fix the `helloocp-cr.yaml` file and apply again: `oc apply -f helloocp-cr.yaml`
- Find the ops metadata.label.app string
- Create a service in console, or on cli:
  - create file `helloocp-service.yaml`:


  - `oc apply -f helloocp-service.yaml`

- create route in console to this service, or on cli by:
  - create file `helloocp-route.yaml`

  - deploy: `oc apply -f helloocp-route.yaml `

- test: `curl http://hello-ocp.apps-crc.testing/hello-ocp`

## 2. Update to reconcile the name to say hello to, from a cr field - required update to model, and reconcile code.

#### Update the model
- Update model in hello_types.go to change the random string to HelloName.
- regen the crd (only changes deploy/crd)
- regen the csv including --update-crds - updates olm-catalog/../manifests/... csv and crd

#### Update the reconcile loop to apply the cr field to a env var in the pod 
- In the Reconcile method, add an env var to the pod spec, set to the cr field helloName
- note: the current Reconcile loop only looks for whether the pod exists. If it does, it doesn't do anything else. We need to ensure we update the pod if the helloName cr field has changed.
  - if rebuild and deploy operator and crd now, you need to delete the pod (remember pods are not reconciled normally, so operator only thing watching it.)
- replace the 'it exists dont requeue' comment and return with an r.client.Update 

- _TODO_ I notice a metrics operator pod/service running - what is this for.

## Creating OLM catsrc, index, bundle
In order to actually make use of OLM, you need to create some or all of these artifacts:

#### OperatorSource
- references a datastore hosting operator bundles
- Creating an operator source automatically creates a catsrc - According to operators book, 
 - not always needed 

#### Catsrc 
- just a yaml file
- references an index image (accord to link: https://docs.openshift.com/container-platform/4.5/operators/olm-managing-custom-catalogs.html#olm-creating-catalog-from-index_olm-managing-custom-catalogs)
- in mq, castro references: ibm-mq-operator-catalog

#### Catalog | index 
- A catalog of operators created from a list of bundles
    - Created with:  opm index add  --bundles quay.io/<namespace>/test-operator:v0.1.0,more…  --tag quay.io/<namespace>/test-catalog:latest 
- (also called index) is created using opm command
- See below for text from ocp doc
- Requires bundle(s) to exist

#### Bundle
- Replaces the deprecated approach of Packagemanifest
- An Operator bundle represents a single version of an Operator. 
- the channel(s) that this version of the Operator are available on are specified as part of creating the bundle
- Bundle manifest includes 1 csv and all csvs crds (so actually the olm/manifests dir)
- Creation with operator-sdk bundle create ….  : https://docs.openshift.com/container-platform/4.5/operators/operator_sdk/osdk-working-bundle-images.html
- Reminder pod man doesn't work on Mac - build uses pod man, MacBook use docker

## Creating them for real

#### bundle
Create the bundle: `operator-sdk bundle create somerandomhostnamem/drb/hello-ocp-operator-bundle:v0.0.1 -b docker`
Push the bundle: `docker push somerandomhostnamem/drb/hello-ocp-operator-bundle:v0.0.1`

##### notes: 
- looks like we add channels by supplying extra args to operator-sdk- bundle create `--channels <channels> --default-channel v0.0.1`
- if not specified, it seems to default to `stable` 

#### catalog
Create the catalog: `opm index add --bundles somerandomhostnamem/drb/hello-ocp-operator-bundle:v0.0.1 --tag somerandomhostnamem/drb/hello-ocp-operator-catalog:v0.0.1 -u docker`

#### catsrc
Enter catalog image into helloocp-catalog-source.yaml

## Trying out OLM operator
- Apply catsrc
- Wait til catsrc ready
- Create operator on stable in specific namespace
- Click on installed operator > Hello OCP
  - observations: 
    - while there is a `+ XCreater instance` button, there is no Hello OCP tab for some reason  **TODO**
    - the sample cr in red ui was not from csv - I believe it is from the csv alm-examples
-   

## reconcile service and route

#### references

- API Ref for 1.17: https://v1-17.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/
- golang route api: https://godoc.org/github.com/openshift/api/route/v1
- new location of golang k8 doc (1.17): https://pkg.go.dev/k8s.io/api@v0.17.2/core/v1

#### watch for service and route

- change the alm-example while I remember, in the csv

- In helloocp_controller.go, add watches for Service and Route types to add method
- Copy existing watch statements, changing the type to the appropriate golang k8/openshift api type (and adding import)
  - imports:
  ```
  corev1 "k8s.io/api/core/v1"       // already exists
  routev1 "route.openshift.io/v1"       // need to add
  ```
  - watch first line:
    - service: `err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{`
      - route: `err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{`

#### Add reconcile logic to Reconcile()

- Move all of pod logic to reconcilePod, ensuring all types carried across as params
- Requeue logic with multiple resources - requeue as soon as hit error, else continue to next resource
- Copy reconcilePod, changing all types to Service, and populating service object with values from helloocp-service.yaml
