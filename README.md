# EnMasse Service Broker

An [Open Service Broker](https://github.com/openservicebrokerapi/servicebroker) implementation.

## Prerequisites

[glide](https://glide.sh/) is used for dependency management. Binaries are available on the
[releases page](https://github.com/Masterminds/glide/releases).


## Setup

```
mkdir -p $GOPATH/src/github.com/EnMasseProject
git clone https://github.com/EnMasseProject/maas-service-broker.git $GOPATH/src/github.com/EnMasseProject/maas-service-broker`
cd $GOPATH/src/github.com/EnMasseProject/maas-service-broker && glide install
```

## Targets

`make run`: Runs the broker with the default profile, configured via `/etc/dev.config.yaml`
`make run-mock-registry`: Mock registry. Entirely separate binary.
`make test`: Runs the test suite.

**Note**

Scripts found in `/test` can act as manual Service Catalog requests until a larger
user scenario can be scripted.

## Setting up

- Run Kubernetes cluster (e.g. through minikube)
- Download and install Helm (https://github.com/kubernetes/helm/releases)
- Run `helm init` to install Helm's Tiller server into the Kubernetes cluster
- Get kubectl, which works with the Service Catalog API server:
  - `docker cp $(docker create duglin/kubectl:latest):kubectl ./servicecatalogctl`
- Configure sc alias and make it connect to the Service Catalog API server:
  - `alias sc="$(pwd)/servicecatalogctl --server=https://$(oc get route apiserver -o jsonpath=\"{.spec.host}\") --insecure-skip-tls-verify"`

## Run Service Catalog and MaaS Broker

- Build Service Catalog images and push them to a registry:
  - Use local docker daemon for building service catalog, NOT the one in minikube, because the build process mounts local volumes into the build container!
  - `REGISTRY=docker.io/luksa make push`
- Run the Service Catalog components:
  - `helm install --namespace=service-catalog --set registry=docker.io/luksa,version=a8d60d7-dirty,storageType=etcd,insecure=true,debug=true,insecureServicePort=31111 deploy/wip-catalog`
- Build MaaS Broker Image & push to registry
  - `make push`
- Deploy MaaS Broker
  - `kubectl create -f kubernetes-resources/maas-broker-deployment.yaml`
- Register Broker in the Service Catalog:
  - `sc create -f broker.yaml`

## Create queue service instance

- `sc create -f instance.yaml`



## Running in OpenShift

```
oc login
oc create -f kubernetes-resources/servicecatalog.yaml

```