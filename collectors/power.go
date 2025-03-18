package collectors

import "github.com/prometheus/client_golang/prometheus"

type PowerCollector struct {
	PowerSupplies map[string]PowerSupply `json:"powerSupplies"`
}

type PowerSupply struct {
	OutputPower   float64                 `json:"outputPower"`
	State         string                  `json:"state"`
	ModelName     string                  `json:"modelName"`
	Capacity      int                     `json:"capacity"`
	InputCurrent  float64                 `json:"inputCurrent"`
	OutputCurrent float64                 `json:"outputCurrent"`
	Uptime        string                  `json:"uptime"`
	Managed       bool                    `json:"managed"`
	TempSensors   map[string]TempSensor   `json:"tempSensors"`
	Fans          map[string]PSUFanStatus `json:"fans"`
}

type PSUFanStatus struct {
	Status string `json:"status"`
	Speed  int    `json:"speed"`
}

type TempSensor struct {
	Status      string `json:"status"`
	Temperature int    `json:"temperature"`
}

func (c *PowerCollector) GetCmd() string {
	return "show environment power"
}

func (c *PowerCollector) Register(registry *prometheus.Registry) {

}

func (c *PowerCollector) UpdateMetrics() {

}
