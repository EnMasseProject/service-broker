# Example payloads for the MaaS Service Broker

These can be posted to the service broker REST API.

To provision a queue:

`curl -H "X-Broker-API-Version: 2.11" -X PUT -H "content-type: application/json" --data-binary @provision-vanilla-queue.json http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca1`

To deprovision the queue:

`curl -H "X-Broker-API-Version: 2.9" -X DELETE "http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca1?service_id=7739ea7d-8de4-4fe8-8297-90f703904589&plan_id=83fc2eaf-d968-4f7d-bbcd-da697ca9232c"`

