package maas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"github.com/pborman/uuid"
	"io"
	"net/http"
)

type MaasClientConfig struct {
	Url string
}

type MaasClient struct {
	config MaasClientConfig
	log    *logging.Logger
}

func NewMaasClient(config MaasClientConfig, log *logging.Logger) (*MaasClient, error) {
	client := &MaasClient{
		config: config,
		log:    log,
	}

	log.Notice("MaaS API Server is at %s", config.Url)

	return client, nil
}

func (c *MaasClient) GetFlavors() ([]Flavor, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v3/flavor", c.config.Url))
	if err != nil {
		return []Flavor{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var flavorList FlavorList

	err = decoder.Decode(&flavorList)
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	return flavorList.Items, nil
}

func (c *MaasClient) GetAddresses(infraID string) ([]Address, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v3/instance/%s/address", c.config.Url, infraID))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var addressList AddressList
	err = decoder.Decode(&addressList)
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	c.log.Infof("Received addresses: %+v", addressList)

	return addressList.Items, nil
}

func (c *MaasClient) GetInstances() ([]Instance, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v3/instance", c.config.Url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var instanceList InstanceList
	err = decoder.Decode(&instanceList)
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	c.log.Infof("Received instances: %+v", instanceList)

	return instanceList.Items, nil
}

func (c *MaasClient) GetInstance(id string) (*Instance, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v3/instance/%s", c.config.Url, id))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	decoder := json.NewDecoder(resp.Body)
	var instance Instance
	err = decoder.Decode(&instance)
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	c.log.Infof("Got instance: %+v", instance)

	return &instance, nil
}

func (c *MaasClient) ProvisionMaaSInfra(infraID string) error {
	c.log.Infof("Provisioning MaaS infrastructure instance %s", infraID)

	instance := Instance{
		Metadata: Metadata{
			Name: infraID,
		},
		Spec: InstanceSpec{
			Namespace: "enmasse-" + infraID,
		},
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(instance)

	c.log.Infof("Sending request: %+v", b)
	resp, err := http.Post(fmt.Sprintf("%s/v3/instance", c.config.Url), "application/json", b)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	c.log.Info("Request sent")

	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//s := buf.String()
	//
	//c.log.Infof("Received response: %s", s)

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&instance)
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	c.log.Infof("Received instance: %+v", instance)

	return nil
}

func (c *MaasClient) ProvisionAddress(infraID string, instanceUUID uuid.UUID, name string, storeAndForward bool, multicast bool, flavor string) error {
	c.log.Infof("Provisioning address %s of flavor %s (instance UUID: %s)", name, flavor, instanceUUID)

	queue := Address{
		Metadata: Metadata{
			Name: name,
			Uuid: instanceUUID.String(),
		},
		Spec: AddressSpec{
			StoreAndForward: storeAndForward,
			Multicast:       multicast,
			Flavor:          flavor,
		},
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(queue)

	c.log.Infof("Sending request: %+v", b)
	resp, err := http.Post(fmt.Sprintf("%s/v3/instance/%s/address", c.config.Url, infraID), "application/json", b)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	c.log.Info("Request sent")

	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//s := buf.String()
	//
	//c.log.Infof("Received response: %s", s)

	decoder := json.NewDecoder(resp.Body)

	var addresses AddressList
	err = decoder.Decode(&addresses)
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	c.log.Infof("Addresses after provisioning: %+v", addresses)

	return nil
}

func (c *MaasClient) DeprovisionAddress(infraID string, instanceUUID uuid.UUID) error {
	c.log.Infof("Deprovisioning address %s", instanceUUID)
	address, err := c.GetAddress(infraID, instanceUUID)
	c.log.Infof("Address name is %s (UUID is %s)", address.Metadata.Name, address.Metadata.Uuid)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v3/instance/%s/address/%s", c.config.Url, infraID, address.Metadata.Name), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	c.log.Infof("Received response: %+v", resp)

	return nil
}

func (c *MaasClient) GetFlavor(planUUID uuid.UUID) (*Flavor, error) {
	flavors, err := c.GetFlavors()
	if err != nil {
		return nil, err
	}
	for _, flavor := range flavors {
		if flavor.Metadata.Uuid == planUUID.String() {
			return &flavor, nil
		}
	}
	return nil, nil
}

// TODO: replace this with more efficient mechanism for looking up addresses across instances
func (c *MaasClient) FindAddress(instanceUUID uuid.UUID) (*Instance, *Address, error) {
	instances, err := c.GetInstances()
	if err != nil {
		return nil, nil, err
	}

	for _, instance := range instances {
		infraID := instance.Metadata.Name
		address, err := c.GetAddress(infraID, instanceUUID)
		if err != nil {
			return nil, nil, err
		}
		if address != nil {
			return &instance, address, nil
		}
	}
	return nil, nil, nil
}

func (c *MaasClient) GetAddress(infraID string, instanceUUID uuid.UUID) (*Address, error) {
	addresses, err := c.GetAddresses(infraID)
	if err != nil {
		return nil, err
	}
	for _, address := range addresses {
		if address.Metadata.Uuid == instanceUUID.String() {
			return &address, nil
		}
	}
	return nil, nil

}
