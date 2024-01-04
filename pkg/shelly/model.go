package shelly

type Status struct {
	MAC string

	Cloud CloudStatus
	MQTT  MQTTStatus
	Bat   BatteryStatus

	HasUpdate bool `json:"has_update"`

	Meters      []Meter
	Relays      []Relay
	Thermostats []Thermostats

	RAMTotal int `json:"ram_total"`
	RAMFree  int `json:"ram_free"`
	FSSize   int `json:"fs_size"`
	FSFree   int `json:"fs_free"`
	Uptime   int
}

type Meter struct {
	Power float64
	Total float64
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

type BatteryStatus struct {
	Value int
}

type Thermostats struct {
	ThermostatPosition float64 `json:"pos"`

	Target_t struct {
		TargetEnabled bool    `json:"enabled"`
		TargetValue   float64 `json:"value"`
		TargetUnit    string  `json:"units"`
	}

	Tmp struct {
		TemperatureValue float64 `json:"value"`
		TemperatureUnit  string  `json:"units"`
	}
}
