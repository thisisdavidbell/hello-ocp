# 1 - Creating the operator

## 0. Prereqs

- a dockerised go application, with image pushed to an accessible docker registry. You can complete chapter 1 - hello go app here: [README-1-hello-go-app.md](README-1-hello-go-app.md) to achieve this, or just use the completed example here: [hello-go-app](hello-go-app).
- an OCP 4.4 or above cluster.
  - The tutorial has been tested on a range of OCP providers, and at OCP 4.4, 4.5, 4.6.
- docker running locally
- operator-sdk installed - https://sdk.operatorframework.io/docs/installation/install-operator-sdk/. The latest version of this tutorial was created with operator sdk 1.0. 
- You must overwrite key values in the examples
  - This tutorial requires a number of values which are specific to your OCP system, such as hostname. Therefore, in all chapters of this tutorial, you should search and replace the following values appropriately:
    - `DOCKERHOSTNAME` - the fully qualified hostname of your docker repo
    - `DOCKERNAMESPACE` - the namespace within your docker repo where all images will be pushed.
    - `OCPHOSTNAME` - the full comain name of your ocp system, specifically everything after `api.` or `apps.` in the urls you are using.)

## Useful resources:
It is highly recommended you read through at least one example of creating an operator. The etcd examples in either of the following 2 sources are a great start to understanding operators. This tutorial assumes you have aleady build up some level of knowledge of kubernetes or OCP, and operators, and so have some understanding of terms such as pod, service, route, operator, operator-sdk, reconcile loop, bundle image, catalog image, catalog source.

- Kubernetes Operators O'Reilly book by Jason Dobies & Joshua Wood - https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/ 
  - this is a fantastic introduction to operators, and will provide all the background knowledge you need. Note however that it uses the deprecated package manifest approach to packaging rather than the newer recommended bundle approach. We will use the bundle approach in this tutorial.
- https://operatorframework.io/ - A great source of help and tutorials
  - https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/ this is the latest version of the more in depth go operator tutorial, creating a basic operator for the existing etcd container.

## Useful notes:
- OLM is installed by default on OCP 4.4 and higher, so you do not need to manually install this as instructed in some of the links above.

## 1. Create an operator to manage the hello go application


## Create a new empty directory
Create a new empty directory inside your $GOPATH/src directory structure, in which you will create all the artifacts in this chapter. If you wish to skip actually creating the code and artifacts in this chapter, you can view or use the finished items in [checkpoints/2-create-the-operator](checkpoints/2-create-the-operator).

Ensure go module support is on by setting:
`export GO111MODULE=on`

### Initialise the project
Create your initial operator project
`operator-sdk init`
note: this used the command `new` in earlier versions of operator-sdk`

### Create your resource definition and controller
You need to define a new resource (a Custom Resource Definition (CRD) describing the custom resource - and so the yaml the user will be able to specify when creating a hello application using your operator), and a controller (the operator code which manages the hello application)

`operator-sdk create api --resource --controller --version=v1alpha1 --group=hellogroup --kind=Hello`
notes: 
  - this used the command `add` in earlier versions of operator-sdk`
  - group is now required
  - `--controller` to create the controller.
  - `--resource` to create the crd

The most important files generated are:
- `api/v1alpha1/hello_types.go` - this defines the various objects that make up the model/api, in go.
- `controllers/hello_controller.go` - the primary go code making up the operator, namely the reconcile methods.

### View the model
- open `api/VERSION/hello_types.go`
- Note the `Hello` structure contains `HelloSpec` and `HelloStatus` child structures. We will change both of these in future chapters, but leave them unchanged for now.
  - `HelloSpec` - this defines the fields of the hello custom resource (cr) that the user can provide to define their hello instance.
  - `HelloStatus` - this defines the fields of the hello custom resource that the operator can use to provide status information on the hello instance
  
### Update model based generated code
- whenever the model/api (for now just 1 file - `api/VERSION/hello_types.go`) is updated, you must run the following to update the generated code:
  - `make generate`
- Similarly, following changes you should always update the crd by running:
  - `make manifests`
    - (in operator-sdk v1.0.1) the crd is defined in `config/crd/bases/hellogroup.my.domain_hellos.yaml`

### Write the reconcile logic
Note: the basic approach described here in this tutorial is based on the tutorial in the [Useful resources](README.md#useful-resources) section above. This links out to a useful example controller here: https://github.com/operator-framework/operator-sdk/blob/master/example/memcached-operator/memcached_controller.go.tmpl if you wish to look at a more detailed example at this point.

- view the file `controllers/hello_controller.go`
- Note it contains no real reconcile logic, just
```
	// your logic here
```
- Write logic to reconcile the Hello kind. The recommended flow is as follows:
  - reconcile the hello object
    - note: in this case we simply want to get the object so that we have the cr yaml. IF it doesnt exist, it must mean its been deliberately deleted. 
  - reconcile the primary object (for example a pod or deployment.) We will use a deployment here.
    - check if the deployment exists 
      - for now, if it does exist, we do nothing else. In later chapters we will make changes to the deployment based on values specified by the user in the hello cr yaml.
      - create primary object (the deployment) if it doesnt exist

Code snippets are shown below for each section. You can also view the finished code for this chapter in [checkpoints/2-create-the-operator](checkpoints/2-create-the-operator).

- reconcile the hello object
```
// Fetch the Hello instance
instance := &hellogroupv1alpha1.Hello{}
err := r.Get(context.TODO(), req.NamespacedName, instance)
if err != nil {
  if errors.IsNotFound(err) {
    // Request object not found, could have been deleted after reconcile request.
    // Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
    // Return and don't requeue
    return ctrl.Result{}, nil
  }
  // Error reading the object - requeue the request.
  return ctrl.Result{}, err
}
```

- reconcile the deployment
Note: to make our life easier later on when we have to reconcile other objects too, we will reconcile the deployment in its own method/
- call out the the new method in the reconcile loop after the above code to reconcile (and so retrieve) the hello instance
```
// Check if the deployment already exists, if not create a new one
result, err := r.reconcileDeployment(instance)
if err != nil || result.Requeue == true {
  // reconcile requested requeue or errored, so requeue
  return result, err
} // else all good carry on
```

- create the new method
  - check if the deployment exists (for the specified name and namespace)


```
func (r *HelloReconciler) reconcileDeployment(instance *hellogroupv1alpha1.Hello) (ctrl.Result, error) {

  found := &appsv1.Deployment{}
  err := r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
  if err != nil && errors.IsNotFound(err) {
    // deployment does not exist 
    //  - we will need to create it shortly and insert the code here
  } else if err != nil {
    // Failed to get Deployment - we got an unexpected error back
    return ctrl.Result{}, err
  }
  // deployment already exists - no error returned, so the get succeeded
  // for now we are happy and dont make any changes to the existing deployments
  return ctrl.Result{}, nil
}
```

- in the case that the deployment doesnt exist, we need to create it. 
  - create a method which creates a deployment for our hello app, which looks how we want a deployment to look based on the hello instance (i.e. the cr yaml)
```
func (r *HelloReconciler) newDeploymentForHello(instance *hellogroupv1alpha1.Hello) *appsv1.Deployment {
	labels := map[string]string{
		"app": instance.Name,
	}
	var replicas int32 = 1

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "hello",
							Image:   "DOCKERHOSTNAME/DOCKERNAMESPACE/hello:v0.0.1",
							Command: []string{"./hello"},
						},
					},
				},
			},
		},
	}
	return dep
}
``` 
  
  - within the deployment reconcile method, call the new method to create a deployment object for the hello application, and set Hello instance as the owner and controller
  
```
// Define a new deployment
dep := r.newDeploymentForHello(instance)

// Set Hello instance as the owner and controller
ctrl.SetControllerReference(instance, dep, r.Scheme)
```

  - within the earlier code, where we found the deployment did not exist, now create the deployment, using our above deployment `dep`. Insert For example, insert the following code:
```
log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
err = r.Create(context.TODO(), dep)
if err != nil {
  log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
  return ctrl.Result{}, err
}
// Deployment created successfully - return and requeue
return ctrl.Result{Requeue: true}, nil
```

_HERE ToDo-next - try out all code to here. - build and push hello app._
_ToDo-next - work out how to follow the latest instructions from tutorial ands my steps below to build and deploy using olm and using operator-sdk. Optionally I could do none-olm install, then olm install?_

_OLD_STEPS_BELOW_HERE_

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
