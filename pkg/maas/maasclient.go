package maas

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/op/go-logging"
	"fmt"
	"bytes"
	"github.com/pborman/uuid"
)

type MaasClientConfig struct {
	Url  string
}

type MaasClient struct {
	config MaasClientConfig
	log      *logging.Logger
}

func NewMaasClient(config MaasClientConfig, log *logging.Logger) (*MaasClient, error) {
	client := &MaasClient{
		config: config,
		log: log,
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

	err = decoder.Decode(&flavorList);
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	flavors := make([]Flavor, len(flavorList.Flavors))
	var i int = 0
	for name := range flavorList.Flavors {
		flavor := flavorList.Flavors[name]
		flavor.Name = name
		//c.log.Infof("Received flavor %v", flavor)
		flavors[i] = flavor
		i++

	}
	return flavors, nil
}


func (c *MaasClient) GetAddresses() ([]Address, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v3/address", c.config.Url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var addressList AddressList
	err = decoder.Decode(&addressList);
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	c.log.Infof("Received addresses: %+v", addressList)

	addresses := make([]Address, len(addressList.Addresses))
	var i int = 0
	for name := range addressList.Addresses {
		addressSpec := addressList.Addresses[name]
		//c.log.Infof("Received addressSpec %v", addressSpec)
		addresses[i] = Address{
			Metadata:Metadata{
				Name: name,
				Uuid: addressSpec.Uuid,
			},
			Spec: addressSpec,
		}
		i++

	}
	return addresses, nil
}


func (c *MaasClient) ProvisionAddress(instanceUUID uuid.UUID, name string, storeAndForward bool, multicast bool, flavor string) (error) {
	c.log.Infof("Provisioning address %s of flavor %s (instance UUID: %s)", name, flavor, instanceUUID)

	queue := Address{
		Metadata: Metadata{
			Name:name,
			Uuid: instanceUUID.String(),
		},
		Spec: AddressSpec{
			StoreAndForward: storeAndForward,
			Multicast: multicast,
			Flavor: flavor,
		},
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(queue)

	c.log.Infof("Sending request: %+v", b)
	resp, err := http.Post(fmt.Sprintf("%s/v3/address", c.config.Url), "application/json", b)
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
	err = decoder.Decode(&addresses);
	if err == io.EOF {
	} else if err != nil {
		c.log.Fatal(err)
	}

	c.log.Infof("Received addresses: %+v", addresses)

	return nil
}

func (c *MaasClient) DeprovisionAddress(instanceUUID uuid.UUID) (error) {
	c.log.Infof("Deprovisioning address %s", instanceUUID)
	address, err := c.GetAddress(instanceUUID)
	c.log.Infof("Address name is %s (UUID is %s)", address.Metadata.Name, address.Metadata.Uuid)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v3/address/%s", c.config.Url, address.Metadata.Name), nil)
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
		if flavor.Uuid == planUUID.String() {
			return &flavor, nil
		}
	}
	return nil, nil
}


func (c *MaasClient) GetAddress(instanceUUID uuid.UUID) (*Address, error) {
	addresses, err := c.GetAddresses()
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

