package collectors

type CoolingCollector struct {
	OverrideFanSpeed           int       `json:"overrideFanSpeed"`
	CoolingMode                string    `json:"coolingMode"`
	ShutdownOnInsufficientFans bool      `json:"shutdownOnInsufficientFans"`
	AmbientTemperature         float64   `json:"ambientTemperature"`
	SystemStatus               string    `json:"systemStatus"`
	AirflowDirection           string    `json:"airflowDirection"`
	PowerSupplySlots           []FanSlot `json:"powerSupplySlots"`
	FanTraySlots               []FanSlot `json:"fanTraySlots"`
}

type FanSlot struct {
	Status string      `json:"status"`
	Speed  int         `json:"speed"`
	Label  string      `json:"label"`
	Fans   []FanStatus `json:"fans"`
}

type FanStatus struct {
	Status                    string  `json:"status"`
	Uptime                    float64 `json:"uptime"`
	MaxSpeed                  int     `json:"maxSpeed"`
	ConfiguredSpeed           int     `json:"configuredSpeed"`
	ActualSpeed               int     `json:"actualSpeed"`
	SpeedStable               bool    `json:"speedStable"`
	LastSpeedStableChangeTime float64 `json:"lastSpeedStableChangeTime"`
}

func (c *CoolingCollector) GetCmd() string {
	return "show environment cooling"
}
