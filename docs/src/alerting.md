# Alerting

The following example alerts use the available metrics from the exporter.

## Listing

To determine if your host/ip is listed, use one of the following.

### Prometheus

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

### Prometheus Operator

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

For more details, see the [Operator Alerting documentation](https://prometheus-operator.dev/docs/user-guides/alerting/).

## Errors

> [!TIP]
> New in `v0.13.0`.

The exporter now returns a per-target `luzilla_rbls_errors` gauge with the
same `(rbl, ip, hostname)` labels as `luzilla_rbls_ips_blacklisted`: `1` if
the most recent check failed, `0` otherwise.

To collapse it to one alert per RBL, aggregate to a scalar:

```yaml
- alert: DnsblRblErrors
  expr: sum by (rbl)(luzilla_rbls_errors) > 0
  for: 15m
  labels:
    severity: warning
  annotations:
    summary: RBL {{ $labels.rbl }} is returning errors
``` 