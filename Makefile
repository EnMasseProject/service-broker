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

run-in-openshift:
	oc process -f kubernetes-resources/servicecatalog.yaml | oc create -f -
	oc create -f kubernetes-resources/maas-broker-deployment.yaml

enmasse-infra:
	curl -o enmasse-deploy.sh https://raw.githubusercontent.com/EnMasseProject/enmasse/master/scripts/enmasse-deploy.sh
	sh enmasse-deploy.sh -p enmasse -u developer -m
	oc new-app -n enmasse --template=enmasse-infra



