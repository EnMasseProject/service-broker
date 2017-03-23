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
  - `oadm policy add-cluster-role-to-user view system:serviceaccount:enmasse:default`
- Deploy MaaS Broker
  - `oc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/kubernetes-resources/maas-broker-deployment.yaml`
- Get special version of kubectl, which knows about Service Catalog resources:
  - `docker cp $(docker create duglin/kubectl:latest):kubectl ./servicecatalogctl`
- Configure sc alias and make it connect to the Service Catalog API server:
  - `alias sc="$(pwd)/servicecatalogctl --server=https://$(oc get route apiserver -o jsonpath=\"{.spec.host}\") --insecure-skip-tls-verify"`
- Register Broker in the Service Catalog:
  - `sc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/examples/service-catalog/broker.yaml`

## Verifying if the broker is registered
- Check the status of the broker:
  - `sc get broker -o yaml`
  - The status.conditions.message should say "Successfully fetched catalog from broker"
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
  - The status.conditions.message should show "The instance was provisioned successfully"
- Verify the MaaS infra pods and the broker pod have been created:
  - Login as admin:
    - `oc login -u admin`
  - List projects:
    - `oc get projects`
    - Find a project named like "_enmasse-63a14329_"
  - List pods in said project:
    - `oc get pods -n enmasse-63a14329`
    - One of the pods should be called "my-vanilla-queue-<something>"

## Provisioning a topic in the same network
- Create the service instance:
  - `sc create -f https://raw.githubusercontent.com/EnMasseProject/service-broker/master/examples/service-catalog/instance-topic.yaml -n my-messaging-project`
    
## Deprovisioning a queue
- Delete the instance object:
  - `sc delete instance my-vanilla-queue -n my-messaging-project`
- Verify the broker pod is terminating:
  - `oc get pods -n enmasse-63a14329`
  


