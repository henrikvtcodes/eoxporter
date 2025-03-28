package main

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/aristanetworks/goeapi"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/henrikvtcodes/eoxporter/collectors"
	"github.com/henrikvtcodes/eoxporter/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const DefaultCollectors = "version,cooling,power,temperature"

var (
	eapiConfigPath           = kingpin.Flag("eapiConf", "Path to Arista eAPI config file. If the flag isn't provided, $EAPI_CONF is checked").Default(os.Getenv("EAPI_CONF")).String()
	defaultCollectorsEnabled = kingpin.Flag("collectors", "Comma-separated list of collectors to enable").Default(DefaultCollectors).String()
	listenAddress            = kingpin.Flag("listen", "Address to listen on for HTTP").Default("localhost:9396").String()
)

func main() {
	// Parse CLI flags
	kingpin.Parse()

	// Handle eAPI Config Path Loading
	eapiAbsoluteConfigPath, err := filepath.Abs(*eapiConfigPath)
	if err != nil {
		util.Logger.Fatal().Err(err).Str("path", *eapiConfigPath).Msgf("Invalid eapiConfig Path: %s", eapiAbsoluteConfigPath)
	}
	util.Logger.Info().Msgf("Loading configuration from %s", eapiAbsoluteConfigPath)
	goeapi.LoadConfig(eapiAbsoluteConfigPath)
	util.Logger.Info().Msgf("Valid Targets: %s", strings.Join(goeapi.Connections(), " "))

	// Handle enabled eAPI Collectors
	defaultCollectors := makeCollectors(strings.Split(*defaultCollectorsEnabled, ","))
	util.Logger.Info().Msgf("Default collectors: %v\n", strings.Join(slices.Collect(maps.Keys(defaultCollectors)), " "))

	// Create the API server
	app := fiber.New()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &util.RequestLogger,
	}))
	app.Get("/metrics", adaptor.HTTPHandlerFunc(MetricsHandler(&defaultCollectors)))

	// Handle HTTP requests
	servErr := app.Listen(*listenAddress)
	if servErr != nil {
		util.Logger.Fatal().Err(servErr).Msg("Error running HTTP server")
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
			util.Logger.Warn().Msgf("Invalid Collector %s", collectorName)
		case "default":
			util.Logger.Info().Msgf("Including Default Collectors")
		case "version":
			if collectorMap["version"] == nil {
				collectorMap["version"] = &collectors.VersionCollector{}
			} else {
				util.Logger.Warn().Msgf("Duplicate collector detected: %s", "version")
			}
		case "cooling":
			if collectorMap["cooling"] == nil {
				collectorMap["cooling"] = &collectors.CoolingCollector{}
			} else {
				util.Logger.Warn().Msgf("Duplicate collector detected: %s", "cooling")

			}
		case "power":
			if collectorMap["power"] == nil {
				collectorMap["power"] = &collectors.PowerCollector{}
			} else {
				util.Logger.Warn().Msgf("Duplicate collector detected: %s", "power")
			}
		case "temperature":
			if collectorMap["temperature"] == nil {
				collectorMap["temperature"] = &collectors.TemperatureCollector{}
			} else {
				util.Logger.Warn().Msgf("Duplicate collector detected: %s", "temperature")

			}
		case "interfaces":
			if collectorMap["interfaces"] == nil {
				collectorMap["interfaces"] = &collectors.InterfacesCollector{}
			} else {
				util.Logger.Warn().Msgf("Duplicate collector detected: %s", "interfaces")
			}
		}

	}
	return collectorMap
}

func mergeCollectors(cMap1 map[string]Collector, cMap2 map[string]Collector) map[string]Collector {
	cMapMerged := make(map[string]Collector)
	for name, coll := range cMap1 {
		cMapMerged[name] = coll
	}
	for name, coll := range cMap2 {
		cMapMerged[name] = coll
	}
	return cMapMerged
}

func MetricsHandler(defaultCollectors *map[string]Collector) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		collectorsParam := params.Get("collectors")
		collectorNames := strings.Split(collectorsParam, ",")
		if len(collectorNames) > 0 && collectorNames[0] != "" {
			if strings.Contains(collectorsParam, "default") {
				util.Logger.Info().Msgf("Merging default collectors with target param")
				collectorMap = mergeCollectors(collectorMap, makeCollectors(collectorNames))
			} else {
				util.Logger.Info().Msgf("Using non-default collectors")
				collectorMap = makeCollectors(collectorNames)
			}
		}
		util.Logger.Info().Msgf("Collectors enabled: %v\n", strings.Join(slices.Collect(maps.Keys(collectorMap)), " "))

		// Specific metrics registry to handle this request
		aristaRegistry := prometheus.NewRegistry()

		// Register prometheus metrics and eAPI commands
		util.Logger.Info().Msgf("Attempting to register collectors with Prometheus and eAPI")
		for name, coll := range collectorMap {
			coll.Register(aristaRegistry)
			if aErr := eapiHandle.AddCommand(coll); aErr != nil {
				util.Logger.Error().Err(aErr).Msgf("Failed to add command for collector %s", name)
				http.Error(w, fmt.Sprintf("Failed to add command for collector %s", name), http.StatusInternalServerError)
				return
			}
		}

		// Get data from switch
		if cErr := eapiHandle.Call(); cErr != nil {
			http.Error(w, "Failed to run Arista eAPI Command", http.StatusInternalServerError)
			util.Logger.Error().Err(cErr).Msgf("Arista eAPI Command call Failed")
			return
		}
		util.Logger.Info().Msg("Arista eAPI Command(s) ran successfully")

		// Update metrics
		for _, coll := range collectorMap {
			coll.UpdateMetrics()
		}
		util.Logger.Info().Msg("Prometheus Metrics Updated")

		// Do the HTTP thing
		metricsHandler := promhttp.HandlerFor(aristaRegistry, promhttp.HandlerOpts{})
		metricsHandler.ServeHTTP(w, r)
	}
}
