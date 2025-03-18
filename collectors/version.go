package collectors

import "github.com/prometheus/client_golang/prometheus"

type VersionCollector struct {
	Uptime           float64 `json:"uptime"`
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	BootupTimestamp  float64 `json:"bootupTimestamp"`
	MemTotal         int     `json:"memTotal"`
	MemFree          int     `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	IsIntlVersion    bool    `json:"isIntlVersion"`
	InternalBuildId  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision"`

	unameInfo *prometheus.GaugeVec
	uptime    *prometheus.Gauge
}

func (c *VersionCollector) GetCmd() string {
	return "show version"
}

func (c *VersionCollector) Register(registry *prometheus.Registry) {

}

func (c *VersionCollector) UpdateMetrics() {

}
