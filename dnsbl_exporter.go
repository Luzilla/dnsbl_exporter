package main

import (
	"os"

	"github.com/Luzilla/dnsbl_exporter/app"

	log "github.com/sirupsen/logrus"
)

// The following are customized during build
var exporterName string = "dnsbl-exporter"
var exporterVersion string
var exporterRev string

func main() {
	dnsbl := app.NewApp(exporterName, exporterVersion)
	dnsbl.Bootstrap()

	err := dnsbl.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
