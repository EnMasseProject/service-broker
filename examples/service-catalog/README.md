# Example payloads for the Service Catalog

These can be posted to the service catalog REST API using a special version of kubectl, which supports Service Catalog resources.

To get the special version of kubectl:

`docker cp $(docker create duglin/kubectl:latest):kubectl ./servicecatalogctl`

Configure `sc` alias and make it connect to the Service Catalog API server:

`alias sc="$(pwd)/servicecatalogctl --server=https://$(oc get route apiserver -o jsonpath=\"{.spec.host}\") --insecure-skip-tls-verify"`

Register the MaaS broker:

`sc create -f broker.yaml`

Create a service instance (to provision a queue):

`sc create -f instance-queue.yaml -n some-namespace`

To delete the service instance (deprovision the queue):

`sc delete instance my-vanilla-queue -n some-namespace`

