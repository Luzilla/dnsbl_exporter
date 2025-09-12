# Alerting

The following example alerts use the available metrics from the exporter.

## Prometheus

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

## Prometheus Operator

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