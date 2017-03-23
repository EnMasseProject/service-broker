package broker

import (
	"github.com/EnMasseProject/maas-service-broker/pkg/maas"
	"github.com/op/go-logging"
	"github.com/pborman/uuid"
	"github.com/EnMasseProject/maas-service-broker/pkg/errors"
	"github.com/kubernetes-incubator/service-catalog/.glide/cache/src/https-k8s.io-kubernetes/pkg/util/strings"
)

type Broker interface {
	Bootstrap() (*BootstrapResponse, error)
	Catalog() (*CatalogResponse, error)
	Provision(uuid.UUID, *ProvisionRequest) (*ProvisionResponse, error)
	Update(uuid.UUID, *UpdateRequest) (*UpdateResponse, error)
	Deprovision(instanceUUID uuid.UUID, serviceId string, planId string) (*DeprovisionResponse, error)
	Bind(uuid.UUID, uuid.UUID, *BindRequest) (*BindResponse, error)
	Unbind(uuid.UUID, uuid.UUID) error
}

type MaasBroker struct {
	log      *logging.Logger
	client   *maas.MaasClient
}


// HACK: needed for deprovisioning; TODO: needs to be replaced with proper storage or removed
var infraIDs map[*uuid.UUID]string = make(map[*uuid.UUID]string)


func NewMaasBroker(
	log *logging.Logger,
	client *maas.MaasClient,
) (*MaasBroker, error) {

	broker := &MaasBroker{
		log:    log,
		client: client,
	}

	return broker, nil
}

func (b MaasBroker) Bootstrap() (*BootstrapResponse, error) {
	b.log.Info("MaaSBroker::Bootstrap")
	return &BootstrapResponse{}, nil
}

const AnycastServiceUUID = "ac6348d6-eeea-43e5-9b97-5ed18da5dcaf"
const MulticastServiceUUID = "7739ea7d-8de4-4fe8-8297-90f703904587"
const QueueServiceUUID = "7739ea7d-8de4-4fe8-8297-90f703904589"
const TopicServiceUUID = "7739ea7d-8de4-4fe8-8297-90f703904590"

func (b MaasBroker) Catalog() (*CatalogResponse, error) {
	b.log.Info("MaaSBroker::Catalog")

	queueService := Service{
		ID:   uuid.Parse(QueueServiceUUID),
		Name: "queue",
		Description: "A messaging queue",
		Bindable: false,
		Plans: []Plan{},
		Metadata: make(map[string]interface{}),
	}

	topicService := Service{
		ID:   uuid.Parse(TopicServiceUUID),
		Name: "topic",
		Description: "A messaging topic",
		Bindable: false,
		Plans: []Plan{},
		Metadata: make(map[string]interface{}),
	}

	flavors, err := b.client.GetFlavors()
	if err != nil {
		b.log.Warning("Could not get flavors from MaaS API server: %s", err.Error())	// TODO: fail here instead of returning any/multicast only?
	}

	b.log.Info("Processing flavors")
	for _,flavor := range flavors {
		b.log.Info("Flavor: %s (%s)", flavor.Metadata.Name, flavor.Spec.Description)
		plan := Plan{
			ID:          uuid.Parse(flavor.Metadata.Uuid),
			Name:        SanitizePlanName(flavor.Metadata.Name),
			Description: flavor.Spec.Description,
			Free:        true,
		}
		if flavor.Spec.Type == maas.Queue {
			queueService.Plans = append(queueService.Plans, plan)
		} else if flavor.Spec.Type == maas.Topic {
			topicService.Plans = append(topicService.Plans, plan)
		} else {
			b.log.Warningf("Unknown flavor type %s", flavor.Spec.Type)
		}
	}

	anycastService := Service{
		ID:          uuid.Parse(AnycastServiceUUID),
		Name:        "direct-anycast-network",
			Description: "A brokerless network for direct anycast messaging",
			Bindable:    false,
			Plans: []Plan{{
			ID:          uuid.Parse("914e9acc-242e-42e3-8995-4ec90d928c2b"),
			Name:        "default",
			Description: "Default plan",
			Free:        true,
		}},
			Metadata: make(map[string]interface{}),
	}

	multicastService := Service{
		ID:          uuid.Parse(MulticastServiceUUID),
		Name:        "direct-multicast-network",
		Description: "A brokerless network for direct multicast messaging",
		Bindable:    false,
		Plans: []Plan{{
			ID:          uuid.Parse("6373d6b9-b701-4636-a5ff-dc5b835c9223"),
			Name:        "default",
			Description: "Default plan",
			Free:        true,
		}},
		Metadata: make(map[string]interface{}),
	}

	services := []Service{
		anycastService,
		multicastService,
	}

	b.log.Info("queueService.Plans: %d", len(queueService.Plans))
	b.log.Info("topicService.Plans: %d", len(topicService.Plans))

	if len(queueService.Plans) > 0 {
		services = append(services, queueService)
	}
	if len(topicService.Plans) > 0 {
		services = append(services, topicService)
	}

	return &CatalogResponse{services}, nil
}

func (b MaasBroker) Provision(instanceUUID uuid.UUID, req *ProvisionRequest) (*ProvisionResponse, error) {
	b.log.Info("Provisioning: %v", req)

	if req.OrganizationID == "" {
		req.OrganizationID = "some-unique-guid"
		req.SpaceID = "some-unique-guid"
	}

	// TODO: this is a temporary hack needed because otherwise the resulting address configmap name is too long
	infraID := strings.ShortenString(req.OrganizationID, 8)

	address, err := b.client.GetAddress(infraID, instanceUUID)
	if err != nil {
		return nil, err
	}
	if address != nil {
		// TODO: verify if parameters are different & return HTTP status Conflict or OK if they haven't
		return &ProvisionResponse{Operation: "successful"}, nil
		//return nil, errors.NewServiceInstanceAlreadyExists(instanceUUID.String())
	}

	name := req.Parameters["name"]
	if name == "" {
		return nil, errors.NewBadRequest("Missing parameter: name")
	}

	switch req.ServiceID.String() {
	case AnycastServiceUUID:
		b.client.ProvisionAddress(infraID, instanceUUID, name, false, false,"")
	case MulticastServiceUUID:
		b.client.ProvisionAddress(infraID, instanceUUID, name, false, true,"")
	case QueueServiceUUID:
		flavor, err := b.client.GetFlavor(req.PlanID)
		if err != nil {
			return nil, err
		} else if flavor == nil || flavor.Spec.Type != maas.Queue {
			return nil, errors.NewBadRequest("Invalid plan ID " + req.PlanID.String())
		}
		b.client.ProvisionAddress(infraID, instanceUUID, name, true, false, flavor.Metadata.Name)
	case TopicServiceUUID:
		flavor, _ := b.client.GetFlavor(req.PlanID)
		if err != nil {
			return nil, err
		} else if flavor == nil || flavor.Spec.Type != maas.Topic {
			return nil, errors.NewBadRequest("Invalid plan ID " + req.PlanID.String())
		}
		b.client.ProvisionAddress(infraID, instanceUUID, name, true, true, flavor.Metadata.Name)
	default:
		return nil, errors.NewBadRequest("Unknown service ID " + req.ServiceID.String())
	}

	infraIDs[&instanceUUID] = infraID

	return &ProvisionResponse{Operation: "successful"}, nil
}

func (b MaasBroker) Deprovision(instanceUUID uuid.UUID, serviceId string, planId string) (*DeprovisionResponse, error) {
	b.log.Info("Deprovisioning %s", instanceUUID.String())


	//infraID := infraIDs[&instanceUUID]
	//if infraID == "" {
	//	return nil, errors.NewServiceInstanceGone(instanceUUID.String())
	//}

	infraID, address, err := b.client.FindAddress(instanceUUID)
	if err != nil {
		return nil, err
	}

	if address == nil {
		return nil, errors.NewServiceInstanceGone(instanceUUID.String())
	}

	b.client.DeprovisionAddress(infraID, instanceUUID)

	infraIDs[&instanceUUID] = ""

	return &DeprovisionResponse{Operation: "successful"}, nil
}

func (b MaasBroker) Bind(instanceUUID uuid.UUID, bindingUUID uuid.UUID, req *BindRequest) (*BindResponse, error) {
	return nil, notImplemented
}

func (b MaasBroker) Unbind(instanceUUID uuid.UUID, bindingUUID uuid.UUID) error {
	return notImplemented
}

func (b MaasBroker) Update(instanceUUID uuid.UUID, req *UpdateRequest) (*UpdateResponse, error) {
	return nil, notImplemented
}
