package main

import (
	"net/http"
	"os"

	"github.com/Luzilla/dnsbl_exporter/app"
	"github.com/Luzilla/dnsbl_exporter/collector"
	"github.com/Luzilla/dnsbl_exporter/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

// The following are customized during build
var exporterName string = "dnsbl-exporter"
var exporterVersion string
var exporterRev string

func createCollector(rbls []string, targets []string, resolver string) *collector.RblCollector {
	collector := collector.NewRblCollector(rbls, targets, resolver)

	return collector
}

func createRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

func main() {
	app := app.NewApp(exporterName, exporterVersion)

	app.Action = func(ctx *cli.Context) error {
		// setup logging
		switch ctx.String("log.output") {
		case "stdout":
			log.SetOutput(os.Stdout)
		case "stderr":
			log.SetOutput(os.Stderr)
		default:
			cli.ShowAppHelp(ctx)
			return cli.NewExitError("We currently support only stdout and stderr: --log.output", 2)
		}
		if ctx.Bool("log.debug") {
			log.SetLevel(log.DebugLevel)
		}

		cfgRbls, err := config.LoadFile(ctx.String("config.rbls"), "rbl")
		if err != nil {
			return err
		}

		cfgTargets, err := config.LoadFile(ctx.String("config.targets"), "targets")
		if err != nil {
			return err
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
				<head><title>DNSBL Exporter</title></head>
				<body>
				<h1>` + exporterName + ` @ ` + exporterVersion + `</h1>
				<p><a href="` + ctx.String("web.telemetry-path") + `">Metrics</a></p>
				<p><a href="https://github.com/Luzilla/dnsbl_exporter">Code on Github</a></p>
				</body>
				</html>`))
		})

		rbls := config.GetRbls(cfgRbls)
		targets := config.GetTargets(cfgTargets)

		registry := createRegistry()

		collector := createCollector(rbls, targets, ctx.String("config.dns-resolver"))
		registry.MustRegister(collector)

		registryExporter := createRegistry()

		if ctx.Bool("web.include-exporter-metrics") {
			log.Infoln("Exposing exporter metrics")

			registryExporter.MustRegister(
				prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
				prometheus.NewGoCollector(),
			)
		}

		handler := promhttp.HandlerFor(
			prometheus.Gatherers{
				registry,
				registryExporter,
			},
			promhttp.HandlerOpts{
				ErrorHandling: promhttp.ContinueOnError,
				Registry:      registry,
			},
		)

		http.Handle(ctx.String("web.telemetry-path"), handler)

		err = http.ListenAndServe(ctx.String("web.listen-address"), nil)
		if err != nil {
			return err
		}

		log.Infoln("Listening on", ctx.String("web.listen-address"))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
