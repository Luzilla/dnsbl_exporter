# Exporter operation modes

Exporter can work in 2 modes:

* classic:
  * targets are configured in `targets.ini`
  * `/metrics` endpoint is used to receive metrics
* prober:
  * targets are configured on the Prometheus side
  * `/prober&target=` is used
  * `targets.ini` should be empty

Each operation mode requires an `rbl.ini`, for example:

```ini
[rbl]
server=ix.dnsbl.manitu.net
```

Prober mode provides additional features over classic:

 1. dynamic configuration of targets on Prometheus side
 1. different query/check interval of queries for different targets (e.g. to work around strict rate limits of the DNSBL)
 1. ~~set different query settings by utilizing probes modules~~ (not yet implemented)

## classic

The following examples configure scraping of metrics from this exporter in classic mode.

Example for a `targets.ini`:

```ini
[targets]
server=smtp.fastmail.com
```

### Prometheus

```yaml
scrape_configs:
  - job_name: 'dnsbl-exporter'
    metrics_path: /metrics
    static_configs:
      - targets: ['127.0.0.1:9211']
```

For more details, see the [Prometheus scrape config documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config).

### Prometheus Operator

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

You can use a `ServiceMonitor` or `PodMonitor`, for more details, see the [Operator ServiceMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#servicemonitor) or [Operator PodMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#podmonitor).

## prober

The following examples configure scraping of metrics from this exporter in prober mode.

### Prometheus

```yaml
scrape_configs:
  - job_name: 'dnsbl-exporter-prober'
    metrics_path: /prober
    params:
      module: [ips]
    static_configs:
      - targets:
        - smtp.fastmail.com
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

### Prometheus Operator

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
        - smtp.fastmail.com
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

> In addition, you can use a `ServiceMonitor` or `PodMonitor` to monitor the `dnsbl-exporter`'s operational metrics, for more details, see the [Operator ServiceMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#servicemonitor) or [Operator PodMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#podmonitor).