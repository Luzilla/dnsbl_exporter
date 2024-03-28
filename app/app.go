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
	"github.com/Luzilla/dnsbl_exporter/internal/resolvconf"
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

const resolvConfFile = "/etc/resolv.conf"

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
			Usage:       "IP address[:port] of the resolver to use, use `system` to use a resolve from " + resolvConfFile,
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
		&cli.BoolFlag{
			Name:  "config.domain-based",
			Usage: "RBLS are domain based blacklists.",
			Value: false,
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
		&cli.StringFlag{
			Name:  "log.format",
			Value: "text",
			Usage: "format, text is logfmt or use json",
			Action: func(cCtx *cli.Context, v string) error {
				if v != "text" && v != "json" {
					return cli.Exit("We currently support only text and json: --log.format", 2)
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
	a.App.Action = func(cCtx *cli.Context) error {
		// setup logging
		handler := &slog.HandlerOptions{}
		var writer io.Writer

		if cCtx.Bool("log.debug") {
			handler.Level = slog.LevelDebug
		}

		switch cCtx.String("log.output") {
		case "stdout":
			writer = os.Stdout
		case "stderr":
			writer = os.Stderr
		}

		var logHandler slog.Handler
		if cCtx.String("log.format") == "text" {
			logHandler = handler.NewTextHandler(writer)
		} else {
			logHandler = handler.NewJSONHandler(writer)
		}

		log := slog.New(logHandler)

		c := config.Config{
			Logger: log.With("area", "config"),
		}

		cfgRbls, err := c.LoadFile(cCtx.String("config.rbls"))
		if err != nil {
			return err
		}

		err = c.ValidateConfig(cfgRbls, "rbl")
		if err != nil {
			return fmt.Errorf("unable to load the rbls from the config: %w", err)
		}

		cfgTargets, err := c.LoadFile(cCtx.String("config.targets"))
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

		// use the system's resolver
		if resolver == "system" {
			log.Info("fetching resolver from " + resolvConfFile)
			servers, err := resolvconf.GetServers(resolvConfFile)
			if err != nil {
				return err
			}
			if len(servers) == 0 {
				return fmt.Errorf("unable to return a server from %s", resolvConfFile)
			}

			// pick the first
			resolver = servers[0]
			log.Info("using resolver: " + resolver)
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

		rblCollector := setup.CreateCollector(rbls, targets, cCtx.Bool("config.domain-based"), dnsUtil, log.With("area", "metrics"))
		registry.MustRegister(rblCollector)

		registryExporter := setup.CreateRegistry()

		if cCtx.Bool("web.include-exporter-metrics") {
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

		http.Handle(cCtx.String("web.telemetry-path"), mHandler.Handler())

		pHandler := prober.ProberHandler{
			DNS:         dnsUtil,
			Rbls:        rbls,
			DomainBased: cCtx.Bool("config.domain-based"),
			Logger:      log.With("area", "prober"),
		}
		http.Handle("/prober", pHandler)

		log.Info("starting exporter",
			slog.String("web.listen-address", cCtx.String("web.listen-address")),
			slog.String("resolver", resolver),
		)
		err = http.ListenAndServe(cCtx.String("web.listen-address"), nil)
		if err != nil {
			return err
		}

		return nil
	}
}

func (a *DNSBLApp) Run(args []string) error {
	return a.App.Run(args)
}
