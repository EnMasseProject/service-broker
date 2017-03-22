build: $(shell find cmd pkg)
	CGO_ENABLED=0 GOOS=linux go build ./cmd/broker

# Will default run to dev profile
run: build vendor
	@${GOPATH}/src/github.com/EnMasseProject/maas-service-broker/scripts/runbroker.sh dev

clean:
	@rm -f ${GOPATH}/bin/broker

vendor:
	@glide install

test: vendor
	go test ./pkg/...

image: build
	mv broker dockerbuild
	cp etc/dev.config.yaml dockerbuild
	docker build -t luksa/maas-broker dockerbuild

push: image
	docker push luksa/maas-broker



catalog:
	curl -H "X-Broker-API-Version: 2.9" localhost:1338/v2/catalog

catalog-noheader:
	curl localhost:1338/v2/catalog

queue:
	curl -H "X-Broker-API-Version: 2.9" -X PUT -H "content-type: application/json" --data-binary @provision-vanilla-queue.json http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca1

deprovision-queue:
	curl -H "X-Broker-API-Version: 2.9" -X DELETE "http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca1?service_id=7739ea7d-8de4-4fe8-8297-90f703904589&plan_id=83fc2eaf-d968-4f7d-bbcd-da697ca9232c"

smallqueue:
	curl -H "X-Broker-API-Version: 2.9" -X PUT -H "content-type: application/json" --data-binary @provision-small-queue.json http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca1

topic:
	curl -H "X-Broker-API-Version: 2.9" -X PUT -H "content-type: application/json" --data-binary @provision-topic.json http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca2

deprovision-topic:
	curl -H "X-Broker-API-Version: 2.9" -X DELETE "http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca2?service_id=7739ea7d-8de4-4fe8-8297-90f703904589&plan_id=83fc2eaf-d968-4f7d-bbcd-da697ca9232c"

anycast:
	curl -H "X-Broker-API-Version: 2.9" -X PUT -H "content-type: application/json" --data-binary @provision-anycast.json http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca3

multicast:
	curl -H "X-Broker-API-Version: 2.9" -X PUT -H "content-type: application/json" --data-binary @provision-multicast.json http://localhost:1338/v2/service_instances/881edff6-30be-43a6-8ca5-8855b8e58ca4


