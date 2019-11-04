package main

import (
	"net/http"
	"os"

	"github.com/luzilla/dnsbl_exporter/collector"
	"github.com/luzilla/dnsbl_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "Print the version information.",
	}

	app := cli.NewApp()
	app.Name = exporterName
	app.Version = exporterVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config.dns-resolver",
			Value: "127.0.0.1",
			Usage: "IP address of the resolver to use.",
		},
		cli.StringFlag{
			Name:  "config.rbls",
			Value: "./rbls.ini",
			Usage: "Configuration file which contains RBLs",
		},
		cli.StringFlag{
			Name:  "config.targets",
			Value: "./targets.ini",
			Usage: "Configuration file which contains the targets to check.",
		},
		cli.StringFlag{
			Name:  "web.listen-address",
			Value: ":9211",
			Usage: "Address to listen on for web interface and telemetry.",
		},
		cli.StringFlag{
			Name:  "web.telemetry-path",
			Value: "/metrics",
			Usage: "Path under which to expose metrics.",
		},
		cli.BoolFlag{
			Name:  "web.include-exporter-metrics",
			Usage: "Include metrics about the exporter itself (promhttp_*, process_*, go_*).",
		},
		cli.BoolFlag{
			Name:  "log.debug",
			Usage: "Enable more output in the logs, otherwise INFO.",
		},
		cli.StringFlag{
			Name:  "log.output",
			Value: "stdout",
			Usage: "Destination of our logs: stdout, stderr",
		},
	}

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

		cfgRbls := config.LoadFile(ctx.String("config.rbls"), "rbl")
		cfgTargets := config.LoadFile(ctx.String("config.targets"), "targets")

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

		log.Infoln("Listening on", ctx.String("web.listen-address"))
		log.Fatal(http.ListenAndServe(ctx.String("web.listen-address"), nil))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
