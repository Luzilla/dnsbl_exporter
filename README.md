# dnsbl-exporter

[![CircleCI](https://circleci.com/gh/Luzilla/dnsbl_exporter.svg?style=svg)](https://circleci.com/gh/Luzilla/dnsbl_exporter) [![pr](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml/badge.svg)](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml) [![Maintainability](https://api.codeclimate.com/v1/badges/31b95e6c679f60e30bea/maintainability)](https://codeclimate.com/github/Luzilla/dnsbl_exporter/maintainability)

This is a server which checks the configured hosts against various DNSBL (sometimes refered to as RBLs).

The idea is to scrape `/metrics` using Prometheus to create graphs, alerts, and so on.

**This is an early release. We accept all kinds of feedback - bug reports, PRs, code, docs, ... :)**

## Using

### Configuration

See `rbls.ini` and `targets.ini` files in this repository.

Otherwise:

```
$ dnsbl-exporter -h
...
--config.dns-resolver value  IP address of the resolver to use. (default: "127.0.0.1")
--config.rbls value          Configuration file which contains RBLs (default: "./rbls.ini")
--config.targets value       Configuration file which contains the targets to check. (default: "./targets.ini")
--web.listen-address value   Address to listen on for web interface and telemetry. (default: ":9211")
--web.telemetry-path value   Path under which to expose metrics. (default: "/metrics")
--log.debug                  Enable more output in the logs, otherwise INFO.
--log.output value           Destination of our logs: stdout, stderr (default: "stdout")
--help, -h                   show help
--version, -V                Print the version information.
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

## Releasing

(This is for myself, since I tend to forget things.)

 1. Create a release on Github
 1. Assemble changelog based on PR merges, etc.
 1. Tag must be `v1.0.0` (semantic versioning, prefixed by `v`)
 1. CircleCI will pick it up and build the binaries
