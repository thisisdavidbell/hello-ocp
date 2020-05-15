# Creating an Operator for our hello-ocp application

Following basic tutorial here (same basic tutorial as previous tutorials): https://github.com/operator-framework/getting-started
 - but trying to apply it to our hello-ocp app
 - using crc - has OLM enabled already

Plan:
 - 1. Create an operator to deploy a hello-ocp image (probably as a pod) - DONE
 - 2. Update operator to also create a service and route as part of a 'helloocp' kind
 - 3. Update to create a deployment, and use size and other fields
 - 4. Consider adding validation, as mentioned in crd file `hello-ocp-operator/pkg/apis/helloocp/v1alpha1/helloocp_types.go`: `	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html`

## 1. Create an operator to deploy a hello-ocp image (probably as a pod)

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

 	// Size is the size of the memcached deployment
 	Size int32 `json:"size"`

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
- `docker login -u kubeadmin -t <oc whoami -t> default-route-openshift-image-registry.apps-crc.testing/project1/hello-ocp`
- `docker push default-route-openshift-image-registry.apps-crc.testing/project1/hello-ocp:v0.0.1`

- register crd: `oc create -f deploy/crds/helloocp.example.com_helloocps_crd.yaml`
- build operator: `operator-sdk build default-route-openshift-image-registry.apps-crc.testing/project1/hello-ocp-operator:v0.0.1`
- push image: `docker push default-route-openshift-image-registry.apps-crc.testing/project1/hello-ocp-operator:v0.0.1`
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

- _TODO_ NEXT STEP - continue tutorial undeploying all objects, them creating clusterserviceversion and deploying operator with OLM.


- _TODO_ I notice a metrics operator pod/service running - what is this for.
