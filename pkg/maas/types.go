package maas

const (
	Queue string = "queue"
	Topic string = "topic"
)

type Flavor struct {
	Name string  `json:"name"`
	Uuid string `json:"uuid"`
	Description string  `json:"description"`
	TemplateName string `json:"templateName"`
	TemplateParameters map[string]string  `json:"templateParameters"`
	Type string  `json:"type"`
}

type FlavorList struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Flavors map[string]Flavor `json:"flavors"`
}


type Metadata struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

type AddressSpec struct {
	StoreAndForward bool `json:"store_and_forward"`
	Multicast bool `json:"multicast"`
	Flavor string `json:"flavor,omitempty"`
	Group string `json:"group,omitempty"`
	Uuid string `json:"uuid,omitempty"`
}

type Address struct {
	Metadata Metadata `json:"metadata"`
	Spec AddressSpec `json:"spec"`
}

type AddressList struct {
	Addresses map[string]AddressSpec `json:"addresses"`
}


