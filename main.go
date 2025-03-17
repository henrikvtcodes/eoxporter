package main

import (
	"fmt"
	"github.com/aristanetworks/goeapi"
	"github.com/henrikvtcodes/eoxporter/collectors"
)

func main() {
	fmt.Printf("Connections:%s\n", goeapi.Connections())
	node, err := goeapi.ConnectTo("hudson")
	if err != nil {
		panic(err)
	}

	eapiVersion := &collectors.VersionCollector{}

	handle, _ := node.GetHandle("json")
	handle.AddCommand(eapiVersion)

	if err := handle.Call(); err != nil {
		panic(err)
	}

	println(eapiVersion.Version)
}
