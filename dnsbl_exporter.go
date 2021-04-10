package main

import (
	"os"

	"github.com/Luzilla/dnsbl_exporter/app"
	"github.com/sirupsen/logrus"
)

// The following are customized during build
var exporterName string = "dnsbl-exporter"
var exporterVersion string
var exporterRev string

// global logger
var log = logrus.New()

func main() {
	dnsbl := app.NewApp(exporterName, exporterVersion, log)
	dnsbl.Bootstrap()

	err := dnsbl.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
