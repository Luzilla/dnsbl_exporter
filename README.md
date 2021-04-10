# dnsbl-exporter

[![pr](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml/badge.svg)](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml) [![Maintainability](https://api.codeclimate.com/v1/badges/31b95e6c679f60e30bea/maintainability)](https://codeclimate.com/github/Luzilla/dnsbl_exporter/maintainability) [![Go Report Card](https://goreportcard.com/badge/github.com/Luzilla/dnsbl_exporter)](https://goreportcard.com/report/github.com/Luzilla/dnsbl_exporter) ![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/Luzilla/dnsbl_exporter?include_prereleases&style=social)

This is a server which checks the configured hosts against various DNSBL (sometimes refered to as RBLs).

The idea is to scrape `/metrics` using Prometheus to create graphs, alerts, and so on.

**This is an early release. We accept all kinds of feedback - bug reports, PRs, code, docs, ... :)**

## Using

### Configuration

See `rbls.ini` and `targets.ini` files in this repository. The files follow the nagios format as this exporter is meant to be a drop-in replacement so you can factor out Nagios, one (simple) step at a time. :-)

Otherwise:

```sh
$ dnsbl-exporter -h
...

GLOBAL OPTIONS:
   --config.dns-resolver value     IP address[:port] of the resolver to use. (default: "127.0.0.1:53") [$DNSBL_EXP_RESOLVER]
   --config.rbls value             Configuration file which contains RBLs (default: "./rbls.ini") [$DNSBL_EXP_RBLS]
   --config.targets value          Configuration file which contains the targets to check. (default: "./targets.ini") [$DNSBL_EXP_TARGETS]
   --web.listen-address value      Address to listen on for web interface and telemetry. (default: ":9211") [$DNSBL_EXP_LISTEN]
   --web.telemetry-path value      Path under which to expose metrics. (default: "/metrics")
   --web.include-exporter-metrics  Include metrics about the exporter itself (promhttp_*, process_*, go_*). (default: false)
   --log.debug                     Enable more output in the logs, otherwise INFO. (default: false)
   --log.output value              Destination of our logs: stdout, stderr (default: "stdout")
   --help, -h                      show help (default: false)
   --version, -v                   print the version (default: false)
```

### Running

 1. Go to [release](https://github.com/Luzilla/dnsbl_exporter/releases) and grab a release for your platform.
 1. Get `rbls.ini` and put it next to the binary.
 1. Get `targets.ini`, and customize. Or use the defaults.
 1. `./dnsbl-exporter`

 Go to http://127.0.0.1:9211/ in your browser.

### Caveat

In order to use this, a _proper_ DNS resolver is needed. Proper means: not Google, not Cloudflare, OpenDNS, etc..
Instead use a resolver like Unbound.

To test on OSX, follow these steps:

```
$ brew install unbound
...
$ sudo unbound -d -vvvv
```
(And leave the Terminal open — there will be ample queries and data for you to see and learn from.)

 Verify Unbound is working and resolution is working:

```
 $ dig +short @127.0.0.1 spamhaus.org
192.42.118.104
```

## License / Author

This code is Apache 2.0 licensed.

For questions, comments or anything else, [please get in touch](https://www.luzilla-capital.com).

## Releasing

(This is for myself, since I tend to forget things.)

 1. `git tag -a x.y.z`
 1. `git push --tags`
 1. GitHub Actions/GoReleaser will build a pretty release
