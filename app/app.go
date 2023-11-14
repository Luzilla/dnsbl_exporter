package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Luzilla/dnsbl_exporter/config"
	"github.com/Luzilla/dnsbl_exporter/internal/index"
	"github.com/Luzilla/dnsbl_exporter/internal/metrics"
	"github.com/Luzilla/dnsbl_exporter/internal/prober"
	"github.com/Luzilla/dnsbl_exporter/internal/setup"
	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"

	x "github.com/miekg/dns"
)

type DNSBLApp struct {
	App *cli.App
}

var (
	appName, appVersion, appPath string
	resolver                     string
)

// NewApp ...
func NewApp(name string, version string) DNSBLApp {
	appName = name
	appVersion = version

	a := cli.NewApp()
	a.Name = appName
	a.Version = appVersion
	a.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config.dns-resolver",
			Value:       "127.0.0.1:53",
			Usage:       "IP address[:port] of the resolver to use.",
			EnvVars:     []string{"DNSBL_EXP_RESOLVER"},
			Destination: &resolver,
		},
		&cli.StringFlag{
			Name:    "config.rbls",
			Value:   "./rbls.ini",
			Usage:   "Configuration file which contains RBLs",
			EnvVars: []string{"DNSBL_EXP_RBLS"},
		},
		&cli.StringFlag{
			Name:    "config.targets",
			Value:   "./targets.ini",
			Usage:   "Configuration file which contains the targets to check.",
			EnvVars: []string{"DNSBL_EXP_TARGETS"},
		},
		&cli.StringFlag{
			Name:    "web.listen-address",
			Value:   ":9211",
			Usage:   "Address to listen on for web interface and telemetry.",
			EnvVars: []string{"DNSBL_EXP_LISTEN"},
		},
		&cli.StringFlag{
			Name:        "web.telemetry-path",
			Value:       "/metrics",
			Usage:       "Path under which to expose metrics.",
			Destination: &appPath,
			Action: func(cCtx *cli.Context, v string) error {
				if !strings.HasPrefix(v, "/") {
					return cli.Exit("Missing / to prefix the path: --web.telemetry-path", 2)
				}
				return nil
			},
		},
		&cli.BoolFlag{
			Name:  "web.include-exporter-metrics",
			Usage: "Include metrics about the exporter itself (promhttp_*, process_*, go_*).",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "log.debug",
			Usage: "Enable more output in the logs, otherwise INFO.",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "log.output",
			Value: "stdout",
			Usage: "Destination of our logs: stdout, stderr",
			Action: func(cCtx *cli.Context, v string) error {
				if v != "stdout" && v != "stderr" {
					return cli.Exit("We currently support only stdout and stderr: --log.output", 2)
				}
				return nil
			},
		},
	}

	return DNSBLApp{
		App: a,
	}
}

func (a *DNSBLApp) Bootstrap() {
	a.App.Action = func(ctx *cli.Context) error {
		// setup logging
		handler := &slog.HandlerOptions{}
		var writer io.Writer

		if ctx.Bool("log.debug") {
			handler.Level = slog.LevelDebug
		}

		switch ctx.String("log.output") {
		case "stdout":
			writer = os.Stdout
		case "stderr":
			writer = os.Stderr
		}

		log := slog.New(handler.NewTextHandler(writer))

		c := config.Config{
			Logger: log.With("area", "config"),
		}

		cfgRbls, err := c.LoadFile(ctx.String("config.rbls"))
		if err != nil {
			return err
		}

		err = c.ValidateConfig(cfgRbls, "rbl")
		if err != nil {
			return fmt.Errorf("unable to load the rbls from the config: %w", err)
		}

		cfgTargets, err := c.LoadFile(ctx.String("config.targets"))
		if err != nil {
			return err
		}

		err = c.ValidateConfig(cfgTargets, "targets")
		if err != nil {
			if !errors.Is(err, config.ErrNoServerEntries) && !errors.Is(err, config.ErrNoSuchSection) {
				return err
			}
			log.Info("starting exporter without targets — check the /prober endpoint or correct the .ini file")
		}

		iHandler := index.IndexHandler{
			Name:    appName,
			Version: appVersion,
			Path:    appPath,
		}

		http.HandleFunc("/", iHandler.Handler)

		rbls := c.GetRbls(cfgRbls)
		targets := c.GetTargets(cfgTargets)

		registry := setup.CreateRegistry()

		dnsUtil, err := dns.New(new(x.Client), resolver, log)
		if err != nil {
			log.Error("failed to initialize dns client")
			return err
		}

		rblCollector := setup.CreateCollector(rbls, targets, dnsUtil, log.With("area", "metrics"))
		registry.MustRegister(rblCollector)

		registryExporter := setup.CreateRegistry()

		if ctx.Bool("web.include-exporter-metrics") {
			log.Info("Exposing exporter metrics")

			registryExporter.MustRegister(
				collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
				collectors.NewGoCollector(),
			)
		}

		mHandler := metrics.MetricsHandler{
			Registry:         registry,
			RegistryExporter: registryExporter,
		}

		http.Handle(ctx.String("web.telemetry-path"), mHandler.Handler())

		pHandler := prober.ProberHandler{
			DNS:    dnsUtil,
			Rbls:   rbls,
			Logger: log.With("area", "prober"),
		}
		http.Handle("/prober", pHandler)

		log.Info("starting exporter",
			slog.String("web.listen-address", ctx.String("web.listen-address")),
			slog.String("resolver", resolver),
		)
		err = http.ListenAndServe(ctx.String("web.listen-address"), nil)
		if err != nil {
			return err
		}

		return nil
	}
}

func (a *DNSBLApp) Run(args []string) error {
	return a.App.Run(args)
}
