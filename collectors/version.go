package collectors

type VersionCollector struct {
	Uptime           float64 `json:"uptime"`
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	BootupTimestamp  int     `json:"bootupTimestamp"`
	MemTotal         int     `json:"memTotal"`
	MemFree          int     `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	IsIntlVersion    bool    `json:"isIntlVersion"`
	InternalBuildId  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision"`
}

func (s *VersionCollector) GetCmd() string {
	return "show version"
}
