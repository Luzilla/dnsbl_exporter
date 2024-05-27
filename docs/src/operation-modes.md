# Exporter operation modes

Exporter can work in 2 modes:

* classic:
  * targets are configured in `targets.ini`
  * `/metrics` endpoint is used to receive metrics
* prober:
  * targets are configured on the Prometheus side
  * `/prober&target=` is used
  * `targets.ini` should be empty

Prober mode provides more advantages over classic mode because of:

1. dynamic configuration of targets on Prometheus side without redeploying/reconfiguring exporter itself.
1. ability to have different interval of queries for different targets, useful when some DNSBL have more strict rate limits then others.
1. ability to set different query settings by utilizing probes modules (not yet implemented).

## classic

The following example configure scraping of metrics from this exporter in classic mode.

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

The following example configure scraping of metrics from this exporter in prober mode.

### Prometheus

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

In addition, you can use a `ServiceMonitor` or `PodMonitor` to monitor the `dnsbl_exporter`'s operational metrics, for more details, see the [Operator ServiceMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#servicemonitor) or [Operator PodMonitor documentation](https://prometheus-operator.dev/docs/operator/design/#podmonitor).