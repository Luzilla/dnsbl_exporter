# dnsbl-exporter - The DNS Block List Exporter

[![pr](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml/badge.svg)](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml) [![Maintainability](https://api.codeclimate.com/v1/badges/31b95e6c679f60e30bea/maintainability)](https://codeclimate.com/github/Luzilla/dnsbl_exporter/maintainability) [![Go Report Card](https://goreportcard.com/badge/github.com/Luzilla/dnsbl_exporter)](https://goreportcard.com/report/github.com/Luzilla/dnsbl_exporter) ![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/Luzilla/dnsbl_exporter?include_prereleases&style=social) [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/luzilla)](https://artifacthub.io/packages/helm/luzilla/dnsbl-exporter)

This is a server (aka Prometheus-compatible exporter) which checks the configured hosts against various DNSBL (DNS Block Lists), sometimes referred to as RBLs.

Should you accept this mission, your task is to scrape `/metrics` using Prometheus to create graphs, alerts, and so on.

**This is (still) pretty early software. But I happily accept all kinds of feedback - bug reports, PRs, code, docs, ... :)**

## Usage

### Configuration

See `rbls.ini` and `targets.ini` files in this repository. The files follow the Nagios format as this exporter is meant to be a drop-in replacement so you can factor out Nagios, one (simple) step at a time. 😊

Otherwise:

```sh
$ dnsbl-exporter -h
...
--config.dns-resolver value  IP address of the resolver to use. (default: "127.0.0.1:53")
--config.rbls value          Configuration file which contains RBLs (default: "./rbls.ini")
--config.targets value       Configuration file which contains the targets to check. (default: "./targets.ini")
--config.domain-based        RBLS are domain instead of IP based blocklists (default: false)
--web.listen-address value   Address to listen on for web interface and telemetry. (default: ":9211")
--web.telemetry-path value   Path under which to expose metrics. (default: "/metrics")
--log.debug                  Enable more output in the logs, otherwise INFO.
--log.output value           Destination of our logs: stdout, stderr (default: "stdout")
--help, -h                   show help
--version, -V                Print the version information.
```

#### System resolver
Resolver can use **first** system resolver from `/etc/resolv.conf` automatically, for this you need set `dns-resolver` to `system`.
Note that exporter should have read permissions to `/etc/resolv.conf` file for this feature to work.

Configure resolver as an argument:
- `--config.dns-resolver=system`

Configure resolver as an environment variable:
- `DNSBL_EXP_RESOLVER=system`

### Deployment options

#### Standalone binary

 1. Go to [release](https://github.com/Luzilla/dnsbl_exporter/releases) and grab a release for your platform.
 1. Get `rbls.ini` and put it next to the binary.
 1. Get `targets.ini`, and customize. Or use the defaults.
 1. `./dnsbl-exporter`

 Go to `http://127.0.0.1:9211/` in your browser.

 As option you can configure exporter to run as systemd service.

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

ADD my-target.ini /target.ini
ADD my-rbls.ini /rbls.ini
```

#### Helm

Additionally, a helm chart is provided to run the exporter on Kubernetes.

To get started quickly, an unbound container is installed into the pod alongside the exporter. This unbound acts as a local DNS server to send queries to. You may turn this off with `unbound.enabled=false` and provide your own resolver (via `config.resolver: an.ip.address:port`).

To configure the chart, copy [`chart/values.yaml`](chart/values.yaml) to `values.local.yaml`; for example, to turn off the included unbound and to supply your own resolver, set your own images and last but not least: supply your own _targets_ and RBLs.

Another useful option for advanced users is to add chart as dependency for your chart which allows you to add inplace resources as part of one release:
```
dependencies:
  - name: dnsbl-exporter
    repository: oci://ghcr.io/luzilla/charts
    version: 0.1.0
```

The sources for the helm chart are in [chart](./chart/), to install it, you can inspect the `Chart.yaml` for the version, check the [helm chart repository](https://github.com/orgs/Luzilla/packages/container/package/charts%2Fdnsbl-exporter) or check out [artifact hub](https://artifacthub.io/packages/helm/luzilla/dnsbl-exporter).

The following command creates a `dnsbl-exporter` release which is installed into a namespace called `my-namespace`:

```sh
helm upgrade --install \
    --namespace my-namespace \
    -f ./chart/values.yaml \
    -f ./values.local.yaml \
    dnsbl-exporter oci://ghcr.io/luzilla/charts/dnsbl-exporter --version 0.1.0
```

### Metrics returned by exporter

The individual configured servers and their status are represented by a **gauge**:

```sh
luzilla_rbls_ips_blacklisted{hostname="mail.gmx.net",ip="212.227.17.168",rbl="ix.dnsbl.manitu.net"} 0
```

This represent the server's hostname and the DNSBL in question. `0` (zero) for unlisted and `1` (one) for listed.
Requests to the DNSBL happen in real-time and are not cached. Take this into account and use accordingly.

If the exporter is configured for DNS based blocklists, the ip label represents the return code of the blocklist.

If you happen to be listed — inspect the exporter's logs as they will contain a reason.

### Exporter operation modes

Exporter can work in 2 modes:
* classic - when targets described in `targets.ini` and `/metrics` endpoint is used to recieve metrics for all targets.
* prober - when targets described on Prometheus side and `/prober&target=` is used. To use this mode `targets.ini` should be empty.

Prober mode provides more advantages over classic mode because of:
1. dynamic configuration of targets on Prometheus side without redeploying/reconfiguring exporter itself.
1. ability to have different interval of queries for different targets, useful when some DNSBL have more strict rate limits then others.
1. ability to set different query settings by utilizing probes modules (not yet implemented).

#### classic

The following example configure scraping of metrics from this exporter in classic mode.

##### Prometheus

```yaml
scrape_configs:
  - job_name: 'dnsbl-exporter'
    metrics_path: /metrics
    static_configs:
      - targets: ['127.0.0.1:9211']
```

For more details, see the [Prometheus scrape config documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config).

##### Prometheus Operator

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: dnsbl-exporter
  namespace: dnsbl-exporter
spec:
  endpoints:
    - interval: 30s
      port: http-9211
      scrapeTimeout: 5s
  jobLabel: dnsbl-exporter
  namespaceSelector:
    matchNames:
      - dnsbl-exporter
  selector:
    matchLabels:
      app.kubernetes.io/instance: dnsbl-exporter
      app.kubernetes.io/name: dnsbl-exporter
```

You can use ServiceMonitor or PodMonitor, for more details, see the [Operator ServiceMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#servicemonitor) or [Operator PodMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#podmonitor).

#### prober

The following example configure scraping of metrics from this exporter in prober mode.

##### Prometheus

```yaml
scrape_configs:
  - job_name: 'dnsbl-exporter-prober'
    metrics_path: /probe
    params:
      module: [ips]
    static_configs:
      - targets:
        - 192.0.2.1
        - 192.0.2.2
        - 192.0.2.3
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9211 # The dnsbl exporter's real hostname:port.
  - job_name: 'dnsbl-exporter-metrics' # collect dnsbl exporter's operational metrics.
    static_configs:
      - targets: ['127.0.0.1:9211']
```

For more details, see the [Prometheus scrape config documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config).

##### Prometheus Operator

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
metadata:
  name: dnsbl-exporter-prober
  namespace: dnsbl-exporter
spec:
  interval: 30s
  jobName: dnsbl-exporter-prober
  module: ips
  prober:
    path: /prober
    scheme: http
    url: dnsbl-exporter.dnsbl-exporter.svc:9211 # Kubernetes dnsbl exporter's service
  scrapeTimeout: 5s
  targets:
    staticConfig:
      static:
        - 192.0.2.1
        - 192.0.2.2
        - 192.0.2.3
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: dnsbl-exporter-metrics # collect dnsbl exporter's operational metrics.
  namespace: dnsbl-exporter
spec:
  endpoints:
    - interval: 30s
      port: http-9211
      scrapeTimeout: 5s
  jobLabel: dnsbl-exporter-metrics
  namespaceSelector:
    matchNames:
      - dnsbl-exporter
  selector:
    matchLabels:
      app.kubernetes.io/instance: dnsbl-exporter
      app.kubernetes.io/name: dnsbl-exporter
```

For more details, see the [Operator Probe documentation](https://prometheus-operator.dev/docs/operator/design/#probe).
You can use ServiceMonitor or PodMonitor to monitor exporter's operation metrics, for more details, see the [Operator ServiceMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#servicemonitor) or [Operator PodMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#podmonitor).

### Alerting

The following example alerts use the scraped metrics from this exporter.

#### prometheus

```yaml
alerts:
  groups:
  - name: dnsbl-exporter
    rules:
    - alert: DnsblRblListed
      expr: luzilla_rbls_ips_blacklisted > 0
      for: 15m
      labels:
        severity: critical
      annotations:
        description: {{ $labels.hostname }} ({{ $labels.ip }}) has been blacklisted in {{ $labels.rbl }} for more than 15 minutes.
        summary: Endpoint {{ $labels.hostname }} is blacklisted
        runbook_url: https://example.org/wiki/runbooks
```

For more details, see the [Prometheus Alerting documentation](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/).

#### Prometheus Operator

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: dnsbl-rules
spec:
  groups:
  - name: dnsbl-exporter
    rules:
      - alert: DnsblRblListed
        expr: luzilla_rbls_ips_blacklisted > 0
        for: 15m
        labels:
          severity: critical
        annotations:
          description: {{ $labels.hostname }} ({{ $labels.ip }}) has been blacklisted in {{ $labels.rbl }} for more than 15 minutes.
          summary: Endpoint {{ $labels.hostname }} is blacklisted
          runbook_url: https://example.org/wiki/runbooks
```

For more details, see the [Operator Alertring documentation](https://prometheus-operator.dev/docs/user-guides/alerting/).

### DNS requirements

In order to use the dnsbl exporter, you need use DNS resolver which speaks with authorative NS of RBLs directly, as each RBL has limitation on amount of requests originating from resolver IP.
Using public resolvers like Google, Cloudflare, OpenDNS, Quad9 etc. will result in facing this limitations, as result no proper monitoring can be done.

Our recomendation is to use resolver like [Unbound](https://github.com/NLnetLabs/unbound) without forwarding, which means you will use root NS.

Exclusion: if you have local copy (mirror) of RBL zone synced over rsync or other channels you can configure local rbldnsd to serve this zone and point Unbound to forward exactly this zone to your's rbldnsd.

To install unbound on OSX, follow these steps:

```sh
$ brew install unbound
...
$ sudo unbound -d -vvvv
```

(And leave the Terminal open — there will be ample queries and data for you to see and learn from.)

An alternative to Homebrew is to use Docker; an example image is provided in this repository, it
contains a working configuration — ymmv.

```sh
docker run -p 53:5353/udp ghcr.io/luzilla/unbound:v0.7.0-rc3
```

Verify Unbound is working and resolution is working:

```sh
 $ dig +short @127.0.0.1 spamhaus.org
192.42.118.104
```

## License / Author

This code is Apache 2.0 licensed.

For questions, comments or anything else, [please get in touch](https://www.luzilla-capital.com).

## Releasing new versions

(This is for myself, since I tend to forget things.)

 1. `git tag -a x.y.z`
 1. `git push --tags`
 1. GitHub Actions/GoReleaser will build a pretty release
