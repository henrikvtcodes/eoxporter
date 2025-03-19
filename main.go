package main

import (
	"flag"
	"fmt"
	"github.com/aristanetworks/goeapi"
	"github.com/henrikvtcodes/eoxporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const DefaultCollectors = "version,cooling,power,temperature"

var (
	eapiConfigPath           = flag.String("eapi-conf", os.Getenv("EAPI_CONF"), "Path to Arista eAPI config file")
	defaultCollectorsEnabled = flag.String("collectors", DefaultCollectors, "Comma-separated list of collectors to enable")
	listenAddress            = flag.String("listen-address", "0.0.0.0:9396", "Address to listen on for HTTP")
)

func main() {
	// Parse CLI flags
	flag.Parse()

	// Handle eAPI Config Path Loading
	eapiAbsoluteConfigPath, err := filepath.Abs(*eapiConfigPath)
	if err != nil {
		println(fmt.Errorf("invalid eapiConfig path: %s", eapiAbsoluteConfigPath))
		panic(err)
	}
	log.Default().Printf("Loading configuration from %s\n", eapiAbsoluteConfigPath)
	goeapi.LoadConfig(eapiAbsoluteConfigPath)
	log.Default().Printf("Valid Targets: %s", strings.Join(goeapi.Connections(), " "))

	// Handle enabled eAPI Collectors
	defaultCollectors := makeCollectors(strings.Split(*defaultCollectorsEnabled, ","))
	log.Default().Printf("Default collectors: %v\n", strings.Join(slices.Collect(maps.Keys(defaultCollectors)), " "))

	// Metrics function
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		MetricsHandler(w, r, &defaultCollectors)
	})

	// Handle HTTP requests
	servErr := http.ListenAndServe(*listenAddress, nil)
	if servErr != nil {
		log.Default().Fatalf("[FATAL] HTTP Server Error: %s", err)
	}
}

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

	// Spit out some logs to the console
	log.Default().Printf("Inbound request for target %v\n", target)
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
	log.Default().Printf("Collectors enabled: %v\n", strings.Join(slices.Collect(maps.Keys(collectorMap)), " "))

	// Specific metrics registry to handle this request
	aristaRegistry := prometheus.NewRegistry()

	// Register prometheus metrics and eAPI commands
	log.Default().Print("Attempting to register collectors")
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
		log.Default().Printf("eAPI Command Failed: %e", cErr)
		return
	}
	log.Default().Println("eAPI command(s) ran successfully")

	// Update metrics
	for _, coll := range collectorMap {
		coll.UpdateMetrics()
	}
	log.Default().Println("Metrics updated")

	// Do the HTTP thing
	metricsHandler := promhttp.HandlerFor(aristaRegistry, promhttp.HandlerOpts{})
	metricsHandler.ServeHTTP(w, r)
}
