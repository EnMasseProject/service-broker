# EnMasse Service Broker

An [Open Service Broker](https://github.com/openservicebrokerapi/servicebroker) implementation for the EnMasse project.

## Setting up
- Run OpenShift cluster (e.g. through minishift)
- Deploy EnMasse global infrastructure:
  - `curl -o enmasse-deploy.sh https://raw.githubusercontent.com/EnMasseProject/enmasse/master/scripts/enmasse-deploy.sh`
  - `bash enmasse-deploy.sh -c "https://openshift-master-url:8443" -u yourusername -p enmasse -m`
- Deploy ServiceCatalog
  - `oc process -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/kubernetes-resources/servicecatalog.yaml | oc create -f -`
- Grant permissions to EnMasse address controller and Service Catalog controller manager
  - SSH to OpenShift master node and perform the following two commands: 
  - `oadm policy add-cluster-role-to-user cluster-admin system:serviceaccount:enmasse:enmasse-service-account`
  - `oadm policy add-cluster-role-to-user edit system:serviceaccount:enmasse:default`
- Deploy MaaS Broker
  - `oc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/kubernetes-resources/maas-broker-deployment.yaml`
- Get kubectl 1.6+ (older versions won't work with the Service Catalog API server):
  - `curl -o kubectl https://storage.googleapis.com/kubernetes-release/release/v1.6.0-beta.4/bin/linux/amd64/kubectl ; chmod +x ./kubectl` (replace linux with darwin if using MacOS)
- Configure sc alias and make it connect to the Service Catalog API server:
  - `alias sc="$(pwd)/kubectl --server=https://$(oc get route apiserver -n enmasse -o jsonpath=\"{.spec.host}\") --insecure-skip-tls-verify"`
- Register Broker in the Service Catalog:
  - `sc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/examples/service-catalog/broker.yaml`

At this point, the Service Catalog will contact the broker and retrieve the list of services the broker is providing. 

## Verifying if the broker is registered
- Check the status of the broker:
  - `sc get broker -o yaml`
  - The `status.conditions.message` should say "Successfully fetched catalog from broker"
- Check if there are four service classes:
  - `sc get serviceclasses`
  - The list should include a "queue" and a "topic" class as well as two "direct-*" classes
  
## Provisioning a queue
- Create a new project/namespace:
  - `oc new-project my-messaging-project`
- Create the service instance:
  - `sc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/examples/service-catalog/instance-queue.yaml -n my-messaging-project`
- Check the service instance's status:
  - `sc get instances -n my-messaging-project -o yaml`
  - The `status.conditions.message` should show "The instance was provisioned successfully"
- Verify the MaaS infra pods and the broker pod have been created:
  - Login as admin:
    - `oc login -u admin`
  - List projects:
    - `oc get projects`
    - Find a project named like "_enmasse-63a14329_"
  - List pods in said project:
    - `oc get pods -n enmasse-63a14329`
    - One of the pods should be called "my-vanilla-queue-<something>"

## Binding the queue
- Create the binding:
  - `sc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/examples/service-catalog/binding-queue.yaml -n my-messaging-project`
- Verify the binding's status:
  - `sc get bindings -n my-messaging-project -o yaml`
  - The `status.conditions.message` property should show "Injected bind result"
- Verify the secret has been created:
  - `oc get secret my-vanilla-queue -o yaml`

## Unbinding the queue
- Delete the binding:
  - `sc delete binding my-vanilla-queue-binding`
- Verify the secret has been deleted:
  - `oc get secrets -n my-messaging-project`

## Provisioning a topic in the same network
- Create the service instance:
  - `sc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/examples/service-catalog/instance-topic.yaml -n my-messaging-project`
    
## Deprovisioning a queue
- Delete the instance object:
  - `sc delete instance my-vanilla-queue -n my-messaging-project`
- Verify the broker pod is terminating:
  - `oc get pods -n enmasse-63a14329`
  


