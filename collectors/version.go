package collectors

type VersionCollector struct {
	Uptime           float64 `json:"uptime"`
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	BootupTimestamp  int64   `json:"bootupTimestamp"`
	MemTotal         int64   `json:"memTotal"`
	MemFree          int64   `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	IsIntlVersion    bool    `json:"isIntlVersion"`
	InternalBuildId  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision"`
}

func (s *VersionCollector) GetCmd() string {
	return "show version"
}
