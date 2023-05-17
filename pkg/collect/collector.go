package collect

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/ianunruh/shelly-exporter/pkg/shelly"
)

func NewCollector(client shelly.Client, log *zap.Logger) *Collector {
	return &Collector{
		Client: client,
		Log:    log,

		MeterPower: prometheus.NewDesc("shelly_meter_power_watts",
			"Current real AC power being drawn, in Watts",
			[]string{"device", "meter"}, nil,
		),
		RelayOn: prometheus.NewDesc("shelly_relay_on",
			"Whether the channel is turned ON or OFF",
			[]string{"device", "relay"}, nil,
		),
		MemoryFree: prometheus.NewDesc("shelly_memory_free_kb",
			"Amount of free memory on device in kilobytes",
			[]string{"device"}, nil,
		),
		MemoryTotal: prometheus.NewDesc("shelly_memory_total_kb",
			"Amount of total memory on device in kilobytes",
			[]string{"device"}, nil,
		),
		FilesystemFree: prometheus.NewDesc("shelly_filesystem_free_kb",
			"Amount of free filesystem space on device in kilobytes",
			[]string{"device"}, nil,
		),
		FilesystemTotal: prometheus.NewDesc("shelly_filesystem_total_kb",
			"Amount of total filesystem space on device in kilobytes",
			[]string{"device"}, nil,
		),
		Uptime: prometheus.NewDesc("shelly_uptime_seconds",
			"Uptime of the device in seconds",
			[]string{"device"}, nil,
		),
		HasUpdate: prometheus.NewDesc("shelly_update_available",
			"Whether the device has a firmware update available",
			[]string{"device"}, nil,
		),
		CloudConnected: prometheus.NewDesc("shelly_cloud_connected",
			"Whether the device is connected to Shelly Cloud",
			[]string{"device"}, nil,
		),
		CloudEnabled: prometheus.NewDesc("shelly_cloud_enabled",
			"Whether the device has Shelly Cloud enabled",
			[]string{"device"}, nil,
		),
		MQTTConnected: prometheus.NewDesc("shelly_mqtt_connected",
			"Whether the device is connected to MQTT",
			[]string{"device"}, nil,
		),
		Latency: prometheus.NewDesc("shelly_api_status_latency_ms",
			"Request latency for the status API in milliseconds",
			[]string{"device"}, nil,
		),
		TemperatureValue: prometheus.NewDesc("shelly_temp",
			"Current temperature",
			[]string{"device", "unit"}, nil,
		),
		ThermostatPosition: prometheus.NewDesc("shelly_thermostat_pos",
			"Indicate to what percentage the thermostat is open",
			[]string{"device"}, nil,
		),
		TargetTemperature: prometheus.NewDesc("shelly_temp_target",
			"Show the target temperature",
			[]string{"device", "unit"}, nil,
		),
		TargetEnabled: prometheus.NewDesc("shelly_temp_target_enabled",
			"Whether the target temperature is set",
			[]string{"device"}, nil,
		),
		BatteryStatus: prometheus.NewDesc("shelly_bat",
			"Show the current battery status in percentage",
			[]string{"device"}, nil,
		),
	}
}

type Collector struct {
	Client shelly.Client
	Log    *zap.Logger

	MeterPower         *prometheus.Desc
	RelayOn            *prometheus.Desc
	MemoryFree         *prometheus.Desc
	MemoryTotal        *prometheus.Desc
	FilesystemFree     *prometheus.Desc
	FilesystemTotal    *prometheus.Desc
	Uptime             *prometheus.Desc
	HasUpdate          *prometheus.Desc
	CloudConnected     *prometheus.Desc
	CloudEnabled       *prometheus.Desc
	MQTTConnected      *prometheus.Desc
	Latency            *prometheus.Desc
	TemperatureValue   *prometheus.Desc
	ThermostatPosition *prometheus.Desc
	TargetTemperature  *prometheus.Desc
	TargetEnabled      *prometheus.Desc
	BatteryStatus      *prometheus.Desc
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.MeterPower
	ch <- c.RelayOn
	ch <- c.MemoryFree
	ch <- c.MemoryTotal
	ch <- c.FilesystemFree
	ch <- c.FilesystemTotal
	ch <- c.Uptime
	ch <- c.HasUpdate
	ch <- c.CloudConnected
	ch <- c.CloudEnabled
	ch <- c.MQTTConnected
	ch <- c.Latency
	ch <- c.TemperatureValue
	ch <- c.ThermostatPosition
	ch <- c.TargetTemperature
	ch <- c.TargetEnabled
	ch <- c.BatteryStatus
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	status, err := c.Client.Status()
	latency := time.Since(start)

	if err != nil {
		c.Log.Error("Failed to call status API",
			zap.Error(err))
		return
	}

	ch <- prometheus.MustNewConstMetric(c.Latency, prometheus.GaugeValue, float64(latency.Milliseconds()),
		status.MAC)

	ch <- prometheus.MustNewConstMetric(c.Uptime, prometheus.GaugeValue, float64(status.Uptime),
		status.MAC)

	ch <- prometheus.MustNewConstMetric(c.MemoryFree, prometheus.GaugeValue, float64(status.RAMFree),
		status.MAC)
	ch <- prometheus.MustNewConstMetric(c.MemoryTotal, prometheus.GaugeValue, float64(status.RAMTotal),
		status.MAC)

	ch <- prometheus.MustNewConstMetric(c.FilesystemFree, prometheus.GaugeValue, float64(status.FSFree),
		status.MAC)
	ch <- prometheus.MustNewConstMetric(c.FilesystemTotal, prometheus.GaugeValue, float64(status.FSSize),
		status.MAC)

	ch <- prometheus.MustNewConstMetric(c.HasUpdate, prometheus.GaugeValue, float64FromBool(status.HasUpdate),
		status.MAC)

	ch <- prometheus.MustNewConstMetric(c.CloudConnected, prometheus.GaugeValue, float64FromBool(status.Cloud.Connected),
		status.MAC)
	ch <- prometheus.MustNewConstMetric(c.CloudEnabled, prometheus.GaugeValue, float64FromBool(status.Cloud.Enabled),
		status.MAC)
	ch <- prometheus.MustNewConstMetric(c.MQTTConnected, prometheus.GaugeValue, float64FromBool(status.MQTT.Connected),
		status.MAC)
	ch <- prometheus.MustNewConstMetric(c.BatteryStatus, prometheus.GaugeValue, float64(status.Bat.Value), status.MAC)

	for i, meter := range status.Meters {
		meterID := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(c.MeterPower, prometheus.GaugeValue, meter.Power,
			status.MAC, meterID)
	}

	for i, relay := range status.Relays {
		relayID := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(c.RelayOn, prometheus.GaugeValue, float64FromBool(relay.IsOn),
			status.MAC, relayID)
	}

	for _, thermostats := range status.Thermostats {
		//temperatureID := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(c.TemperatureValue, prometheus.GaugeValue, thermostats.Tmp.TemperatureValue, status.MAC, thermostats.Tmp.TemperatureUnit)
		ch <- prometheus.MustNewConstMetric(c.ThermostatPosition, prometheus.GaugeValue, thermostats.ThermostatPosition, status.MAC)
		ch <- prometheus.MustNewConstMetric(c.TargetTemperature, prometheus.GaugeValue, thermostats.Target_t.TargetValue, status.MAC, thermostats.Target_t.TargetUnit)
		ch <- prometheus.MustNewConstMetric(c.TargetEnabled, prometheus.GaugeValue, float64FromBool(thermostats.Target_t.TargetEnabled), status.MAC)
	}
}

func float64FromBool(value bool) float64 {
	var out float64
	if value {
		out = 1
	}
	return out
}
