package app

import (
	"github.com/urfave/cli"
)

// NewApp ...
func NewApp(name string, version string) *cli.App {

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "Print the version information.",
	}

	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config.dns-resolver",
			Value: "127.0.0.1:53",
			Usage: "IP address[:port] of the resolver to use.",
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

	return app
}
