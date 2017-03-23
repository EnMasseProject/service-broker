# Example payloads for address controller

These can be posted to the address controller REST API.

To create an instance, do the following:

`curl -X POST -H "content-type: application/json" --data-binary @instance.json ${ADDRESS_CONTROLLER_URL}/v3/instance`

To provision a queue:

`curl -X POST -H "content-type: application/json" --data-binary @vanilla-queue.json ${ADDRESS_CONTROLLER_URL}/v3/instance/my-instance/address`

To deprovision the queue:

`curl -X DELETE ${ADDRESS_CONTROLLER_URL}/v3/instance/my-instance/address/my-vanilla-queue`

