# Operators with OLM

# Creating an Operator tutorial

Following: https://docs.openshift.com/container-platform/4.4/operators/operator_sdk/osdk-getting-started.html#managing-memcached-operator-using-olm_osdk-getting-started

Follow instructions in link above, withj following changes:
 - created the operatorGroup in default and project1 namespaces
 - update `namespace: placeholder` to `namespace: project1` in `deploy/olm-catalog/memcached-operator/manifests/memcached-operator.clusterserviceversion.yaml`
 - command seems to be `oc apply -f deploy/olm-catalog/memcached-operator/manifests/memcached-operator.clusterserviceversion.yaml`
 - `oc get csv` shows successful.
 - however, console shows `Cannot update - Catalog source was removed` under 'Installed Operators'
   - creation of 'instances' works fine, and console shows the operator (under Installed Operators) seems to be deployed correctly.
   - google suggests this error relates to subscriptions, which we don't have. Perhaps expected? _TODO - understand this error._

 - for update section, steps don't work. Try instead:
   - _TODO_ understand the intention of this update process better - it seems to be focused on updating the operator itself. Is this the mechanism to update the application? A new operator version per new application version? Migration of app? How is this triggered?
   - duplicate clusterserviceversion file.
     - update `version:` and `metadata.name:` to 0.0.2
     - add line at end: `replaces: memcached-operator.v0.0.1`
   - note current csv version: `oc get csv`
   - note current state of memcached deployments, e.g. `oc get deployments`
   - note current age of operator deployment: `oc get deployments | grep memcached-operator`
   - deploy, e.g.: `oc apply -f deploy/olm-catalog/memcached-operator/manifests/memcached-operator.clusterserviceversion.0.0.2.yaml`
   - only `oc get csv` shows any visible change, initially:
   ```
   oc get csv
NAME                        DISPLAY              VERSION   REPLACES                    PHASE
memcached-operator.v0.0.1   Memcached Operator   0.0.1                                 Deleting
memcached-operator.v0.0.2   Memcached Operator   0.0.2     memcached-operator.v0.0.1   Succeeded
```
then:
```
oc get csv
NAME                        DISPLAY              VERSION   REPLACES                    PHASE
memcached-operator.v0.0.2   Memcached Operator   0.0.2     memcached-operator.v0.0.1   Succeeded
```


To try (and FAILED) to get a visible change:
 - Tag and push new memcached image, i.e.:
 ```
 docker tag default-route-openshift-image-registry.apps-crc.testing/project1/memcached:1.4.36-alpine default-route-openshift-image-registry.apps-crc.testing/project1/memcached:v0.0.3
 docker push default-route-openshift-image-registry.apps-crc.testing/project1/memcached:v0.0.3
 ```
 - Update version: in following section of the csv yaml to `v0.0.3`
 ```
  spec:
  apiservicedefinitions: {}
  customresourcedefinitions:
    owned:
    - description: Memcached is the Schema for the memcacheds API
      kind: Memcached
      name: memcacheds.cache.example.com
      version: v0.0.3
```
This should tell the operator that the version of memcached image to use has changed.
 - update the `replaces:` to `0.0.3`
 - check csv, and deployments again, including: `oc get deployment memcached-for-drupal -o yaml | grep image:`
 - deploy again: `oc apply -f deploy/olm-catalog/memcached-operator/manifests/memcached-operator.clusterserviceversion.0.0.3.yaml`
 - recheck
 initially:
 ```
 $ oc get csv
NAME                        DISPLAY              VERSION   REPLACES                    PHASE
memcached-operator.v0.0.2   Memcached Operator   0.0.2     memcached-operator.v0.0.1   Replacing
memcached-operator.v0.0.3   Memcached Operator   0.0.3     memcached-operator.v0.0.2   Pending
```
This got itself in a mess. Refer to above _TODO_ before spending any more time with updating.

# Create operator managed by OLM for hello-ocp app

Try creating operators with OLM for hello-ocp app following: [OPERATORS-HELLO-OCP.md](OPERATORS-HELLO-OCP.md)
