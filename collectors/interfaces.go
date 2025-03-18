package collectors

import "github.com/prometheus/client_golang/prometheus"

type InterfacesCollector struct {
	Interfaces map[string]Interface `json:"interfaces"`
}

type Interface struct {
	OutBroadcastPackets int     `json:"outBroadcastPkts"`
	OutUnicastPackets   int     `json:"outUcastPkts"`
	OutMulticastPackets int     `json:"outMulticastPkts"`
	OutDiscards         int     `json:"outDiscards"`
	OutOctets           int     `json:"outOctets"`
	InBroadcastPackets  int     `json:"inBroadcastPkts"`
	InUnicastPackets    int     `json:"inUcastPkts"`
	InMulticastPackets  int     `json:"inMulticastPkts"`
	InDiscards          int     `json:"inDiscards"`
	InOctets            int     `json:"inOctets"`
	LastUpdateTimestamp float64 `json:"lastUpdateTimestamp"`
}

func (c *InterfacesCollector) GetCmd() string {
	// In the context of eAPI, this command seems to return the output of "show interfaces counters incoming"
	// and show interfaces counters outgoing" (and the outgoing command does the same)
	return "show interfaces counters incoming"
}

func (c *InterfacesCollector) Register(registry *prometheus.Registry) {

}

func (c *InterfacesCollector) UpdateMetrics() {

}
