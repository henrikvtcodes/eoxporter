package main

import (
	"fmt"
	"github.com/aristanetworks/goeapi"
	"github.com/henrikvtcodes/eoxporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"slices"
	"strings"
)

type Collector interface {
	GetCmd() string
	Register(*prometheus.Registry)
	UpdateMetrics()
}

func makeCollectors(collectorNames []string) map[string]Collector {
	collectorMap := make(map[string]Collector)
	for _, collectorName := range collectorNames {
		switch strings.ToLower(collectorName) {
		default:
			log.Default().Printf("[WARN] Invalid Collector %s", collectorName)
		case "version":
			if collectorMap["version"] == nil {
				collectorMap["version"] = &collectors.VersionCollector{}
			} else {
				log.Default().Printf("[WARN] Duplicate collector: version")
			}
		case "cooling":
			if collectorMap["cooling"] == nil {
				collectorMap["cooling"] = &collectors.CoolingCollector{}
			} else {
				log.Default().Printf("[WARN] Duplicate collector: cooling")
			}
		case "power":
			if collectorMap["power"] == nil {
				collectorMap["power"] = &collectors.PowerCollector{}
			} else {
				log.Default().Printf("[WARN] Duplicate collector: power")
			}
		case "temperature":
			if collectorMap["temperature"] == nil {
				collectorMap["temperature"] = &collectors.TemperatureCollector{}
			} else {
				log.Default().Printf("[WARN] Duplicate collector: temperature")
			}
		}

	}
	return collectorMap
}

func MetricsHandler(w http.ResponseWriter, r *http.Request, defaultCollectors *map[string]Collector) {
	params := r.URL.Query()

	// Get target and validate that it is valid
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}
	if !(slices.Contains(goeapi.Connections(), target)) {
		http.Error(w, "Target does not exist in config", http.StatusNotFound)
		return
	}

	// Initialize eAPI Handle
	node, err := goeapi.ConnectTo(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to %q", target), http.StatusInternalServerError)
		return
	}
	eapiHandle, err := node.GetHandle("json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to %q", target), http.StatusInternalServerError)
		return
	}

	// Figure out which collectors need to handle this request
	collectorMap := *defaultCollectors
	collectorNames := strings.Split(params.Get("collectors"), ",")
	if len(collectorNames) > 0 {
		collectorMap = makeCollectors(collectorNames)
	}

	aristaRegistry := prometheus.NewRegistry()

	// Register prometheus metrics and eAPI commands
	for name, coll := range collectorMap {
		coll.Register(aristaRegistry)
		if aErr := eapiHandle.AddCommand(coll); aErr != nil {
			http.Error(w, fmt.Sprintf("Failed to add command for collector %s", name), http.StatusInternalServerError)
			return
		}
	}

	// Get data from switch
	if cErr := eapiHandle.Call(); cErr != nil {
		http.Error(w, "Failed to run Arista eAPI Command", http.StatusInternalServerError)
		return
	}

	// Update metrics
	for _, coll := range collectorMap {
		coll.Register(aristaRegistry)
	}

	// Do the HTTP thing
	metricsHandler := promhttp.HandlerFor(aristaRegistry, promhttp.HandlerOpts{})
	metricsHandler.ServeHTTP(w, r)
}
