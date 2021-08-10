package shelly

type Status struct {
	MAC string

	Cloud CloudStatus
	MQTT  MQTTStatus

	HasUpdate bool `json:"has_update"`

	Meters []Meter
	Relays []Relay

	RAMTotal int `json:"ram_total"`
	RAMFree  int `json:"ram_free"`
	FSSize   int `json:"fs_size"`
	FSFree   int `json:"fs_free"`
	Uptime   int
}

type Meter struct {
	Power float64
}

type Relay struct {
	IsOn bool
}

type CloudStatus struct {
	Connected bool
	Enabled   bool
}

type MQTTStatus struct {
	Connected bool
}
