package main

import (
	"flag"
	"fmt"
	"github.com/aristanetworks/goeapi"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const DefaultCollectors = "version,cooling,power,temperature"

var (
	eapiConfigPath           = flag.String("eapi-config", os.Getenv("EAPI_CONF"), "Path to Arista eAPI config file")
	defaultCollectorsEnabled = flag.String("collectors", DefaultCollectors, "Comma-separated list of collectors to enable")
	listenAddress            = flag.String("listen-address", "0.0.0.0:9396", "Address to listen on for HTTP")
)

func main() {
	flag.Parse()

	// Handle eAPI Config Path Loading
	eapiAbsoluteConfigPath, err := filepath.Abs(*eapiConfigPath)
	if err != nil {
		println(fmt.Errorf("invalid eapiConfig path: %s", eapiAbsoluteConfigPath))
		panic(err)
	}
	goeapi.LoadConfig(eapiAbsoluteConfigPath)

	// Handle enabled eAPI Collectors
	defaultCollectors := makeCollectors(strings.Split(*defaultCollectorsEnabled, ","))

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		MetricsHandler(w, r, &defaultCollectors)
	})

	// Handle HTTP requests
	servErr := http.ListenAndServe(*listenAddress, nil)
	if servErr != nil {
		log.Default().Fatalf("[FATAL] HTTP Server Error: %s", err)
	}
}
