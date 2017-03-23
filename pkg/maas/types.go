package maas

const (
	Queue string = "queue"
	Topic string = "topic"
)

type Flavor struct {
	Metadata Metadata `json:"metadata"`
	Spec FlavorSpec `json:"spec"`
}

type FlavorSpec struct {
	Type string  `json:"type"`
	Description string  `json:"description"`
	TemplateName string `json:"templateName"`
	TemplateParameters map[string]string  `json:"templateParameters"`
}

type FlavorList struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Items []Flavor `json:"items"`
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
}

type Address struct {
	Metadata Metadata `json:"metadata"`
	Spec AddressSpec `json:"spec"`
}

type AddressList struct {
	Items []Address `json:"items"`
}


type Instance struct {
	Metadata Metadata `json:"metadata"`
	Spec InstanceSpec `json:"spec"`
}

type InstanceSpec struct {
	Namespace string `json:"namespace"`
	MessagingHost bool `json:"messagingHost"`
	MQTTHost bool `json:"mqttHost"`
	ConsoleHost bool `json:"consoleHost"`
}

type InstanceList struct {
	Items []Instance `json:"items"`
}
