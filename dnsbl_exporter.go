package main

import (
	"fmt"
	"os"

	"github.com/Luzilla/dnsbl_exporter/app"
)

// The following are customized during build
var exporterName string = "dnsbl-exporter"
var exporterVersion string = "dev"

func main() {
	dnsbl := app.NewApp(exporterName, exporterVersion)
	dnsbl.Bootstrap()

	err := dnsbl.Run(os.Args)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}
}
