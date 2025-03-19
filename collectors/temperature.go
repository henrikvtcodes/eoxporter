package collectors

import "github.com/prometheus/client_golang/prometheus"

type TemperatureCollector struct {
	ShutdownOnOverheat bool                `json:"shutdownOnOverheat"`
	SystemStatus       string              `json:"systemStatus"`
	TemperatureSensors []TemperatureSensor `json:"tempSensors"`
	PowerSupplySlots   []PSUSlot           `json:"powerSupplySlots"`
}

type PSUSlot struct {
	ENTPhysicalClass   string              `json:"entPhysicalClass"`
	RelativePosition   string              `json:"relPos"`
	TemperatureSensors []TemperatureSensor `json:"tempSensors"`
}

type TemperatureSensor struct {
	MaxTemperature           float64 `json:"maxTemperature"`
	MaxTemperatureLastChange float64 `json:"maxTemperatureLastChange"`
	HwStatus                 string  `json:"hwStatus"`
	AlertCount               int     `json:"alertCount"`
	Description              string  `json:"description"`
	OverheatThreshold        int     `json:"overheatThreshold"`
	CriticalThreshold        int     `json:"criticalThreshold"`
	InAlertState             bool    `json:"inAlertState"`
	TargetTemperature        int     `json:"targetTemperature"`
	RelPos                   string  `json:"relPos"`
	CurrentTemperature       float64 `json:"currentTemperature"`
	PidDriverCount           int     `json:"pidDriverCount"`
	IsPidDriver              bool    `json:"isPidDriver"`
	Name                     string  `json:"name"`
}

func (c *TemperatureCollector) GetCmd() string {
	return "show environment temperature"
}

func (c *TemperatureCollector) Register(registry *prometheus.Registry) {

}

func (c *TemperatureCollector) UpdateMetrics() {

}
