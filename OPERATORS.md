Creating Operators

# Creating an Operator tutorial

Following: https://docs.openshift.com/container-platform/4.4/operators/operator_sdk/osdk-getting-started.html#osdk-getting-started

## 0. Prereqs

  - versions - I found the generation steps only worked when I ensured I used:
    - operator-sdk version = 0.17.0 (from brew install)
    - go version = operator-sdk version = 1.14.2
  - _TODO_ - quay.io not working with crc. Establish how to make quay.io available from within crc - Need a https://quay.io/ (or alternative public image repo)
    - _TODO_ is a quay.io account this something we should create?
    - for now use crc's internal registry

## 1. Create Operator using operator-sdk

Follow tutorial steps at: https://docs.openshift.com/container-platform/4.4/operators/operator_sdk/osdk-getting-started.html#osdk-getting-started

  - _TODO_ - not clear how to try the custom validation
    - assumption that `operator-sdk generate openapi` is now `operator-sdk generate crds`
    - changed `oc create -f deploy/crds/cache_v1alpha1_memcached_crd.yaml` to `oc create -f deploy/crds/cache.example.com_memcacheds_crd.yaml`

  - Attempt option to run operator into cluster
  - use working container registry. For now this may be the crc internal registry not the quay.io one. So pick correct next steps:

  - _TODO - make this option work. AVOID for now_ quay.io
    - when running the `operator-sdk build quay.io/davidrbell/memcached-operator:v0.0.1`, ensure the namespace is your quay.io username
    - manually perform image name update of deploy/operator.yaml
    - run `docker login quay.io` before the push
  - internal registry
    - expose registry on default route:
      - `oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge`
    - grant permission for kubeadmin to use registry:
      - `oc policy add-role-to-user registry-viewer kubeadmin`
    - docker login:
      - `docker login -u kubeadmin -p <output from 'oc whoami -t'> default-route-openshift-image-registry.apps-crc.testing`
      - with crc not having external access, also need to get the memcached image which the memcached operator uses for the memcached containers:
        - `docker pull memcached:1.4.36-alpine`
        - `docker tar memcached:1.4.36-alpine default-route-openshift-image-registry.apps-crc.testing/project1/memcached:1.4.36-alpine`
        - `docker push default-route-openshift-image-registry.apps-crc.testing/project1/memcached:1.4.36-alpine`
        - view image in console, and note internal url:
          - `crc console`
          - Builds -> Image streams -> memcached-operator
          - see 'Image Repository' value
    - update image referenced by operator controller to internal image value, e.g. `image-registry.openshift-image-registry.svc:5000/project1/memcached-operator:v0.0.1`
      - in pkg/controllermemcached/memcached_controller.go
      - update Image line to internal image, e.g.: `default-route-openshift-image-registry.apps-crc.testing/project1/memcached:1.4.36-alpine`
    - Build operator image:
      - `operator-sdk build default-route-openshift-image-registry.apps-crc.testing/project1/memcached-operator:v0.0.1`
    - `docker images | grep memcached`
    - `docker push default-route-openshift-image-registry.apps-crc.testing/project1/memcached-operator:v0.0.1`
    - view image in console, and note internal url:
      - `crc console`
      - Builds -> Image streams -> memcached-operator
      - see 'Image Repository' value
    - manually perform image name update of deploy/operator.yaml - this is the operator image above (not the memcached image!)

  - verify you can complete the tutorial and deploy a memcached-operator deployment
  - note that you can run: `oc get Memcached`
    - I note this does not have the READY UP-TO-DATE AVAILABLE columns that `oc get dpeloyments` has.
  - note in console, 'Installed Operators' is empty - presume this relates to using OLM rather than approach used here. s

## 2. Manage opertaor using Operatopr Lifecycle Manager (OLM)

 - complete the tutorial at: https://docs.openshift.com/container-platform/4.4/operators/operator_sdk/osdk-getting-started.html#managing-memcached-operator-using-olm_osdk-getting-started
