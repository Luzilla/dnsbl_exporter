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
--config.dns-resolver value  IP address of the resolver to use. (default: "127.0.0.1:53")
--config.rbls value          Configuration file which contains RBLs (default: "./rbls.ini")
--config.targets value       Configuration file which contains the targets to check. (default: "./targets.ini")
--config.domain-based        RBLS are domain instead of IP based blacklists (default: false)
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

#### Container

Docker/OCI images are available in the [container registry](https://github.com/orgs/Luzilla/packages?repo_name=dnsbl_exporter):

```sh
$ docker pull ghcr.io/luzilla/dnsbl_exporter:vX.Y.Z
...
```

Please note: `latest` is not provided.

The images expect `target.ini` and `rbls.ini` in the following location:

```sh
/
```

Either start the container and supply the contents, or build your own image:

```sh
docker run \
    --rm \
    -e DNSBL_EXP_RESOLVER=your.resolver:53 \
    -p 9211:9211 \
    -v ./conf:/etc/dnsbl-exporter \
    ghcr.io/luzilla/dnsbl_exporter:vA.B.C
```

```docker
FROM ghcr.io/luzilla/dnsbl_exporter:vA.B.C

ADD my-target.ini /etc/dnsbl-exporter/target.ini
ADD my-rbls.ini /etc/dnsbl-exporter/rbls.ini
```


#### Querying

The individual configured servers and their status are represented by a **gauge**:

```sh
luzilla_rbls_ips_blacklisted{hostname="mail.gmx.net",ip="212.227.17.168",rbl="ix.dnsbl.manitu.net"} 0
```

This represent the server's hostname and the DNSBL in question. `0` for unlisted and `1` for listed. Requests to the DNSBL happen in real-time and are not cached. Take this into account and use accordingly.

If the exporter is configured for DNS based blacklists, the ip label represents the return code of the blacklist.

### Caveat

In order to use this, a _proper_ DNS resolver is needed. Proper means: not Google, not Cloudflare, OpenDNS, etc..
Instead use a resolver like [Unbound](https://github.com/NLnetLabs/unbound).

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
