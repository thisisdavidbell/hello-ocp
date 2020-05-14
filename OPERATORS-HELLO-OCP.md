# Creating an Operator for our hello-ocp application

Following basic tutorial here (same process as previous tutorials): https://github.com/operator-framework/getting-started
 - but trying to apply it to our hello-ocp app
 - using crc - has OLM enabled already

Plan:
 - 1. create an operator to deploy a hello-ocp image (probably as a pod)
 - 2. update operator to also create a route as part of a 'helloocp' instance
 - 3. update to create a deployment, and use size and other fields


Following steps, with changes:
 - command to create crd: `operator-sdk add api --api-version=helloocp.example.com/v1alpha1 --kind=Helloocp`
 - _TODO_ Consider adding validation, as mentioned in crd file `hello-ocp-operator/pkg/apis/helloocp/v1alpha1/helloocp_types.go`: `	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html`
 - adding an extra spec item to see how its handled.

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
 Note, these I believe will be completely ignored as we dont write any code to use them (size makes no sense for a pod anyway).

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
replacing busybox with with image and name (we will push image to internal registry shortly.) so:
```
Spec: corev1.PodSpec{
  Containers: []corev1.Container{
    {
      Name:    "helloocp",
      Image:   "image-registry.openshift-image-registry.svc:5000",
      Command: []string{"sleep", "3600"},
    },
  },
  ```
