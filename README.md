# dnsbl-exporter - The DNS Block List Exporter

[![pr](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml/badge.svg)](https://github.com/Luzilla/dnsbl_exporter/actions/workflows/pr.yml) [![Maintainability](https://api.codeclimate.com/v1/badges/31b95e6c679f60e30bea/maintainability)](https://codeclimate.com/github/Luzilla/dnsbl_exporter/maintainability) [![Go Report Card](https://goreportcard.com/badge/github.com/Luzilla/dnsbl_exporter)](https://goreportcard.com/report/github.com/Luzilla/dnsbl_exporter) ![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/Luzilla/dnsbl_exporter?include_prereleases&style=social) [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/luzilla)](https://artifacthub.io/packages/helm/luzilla/dnsbl-exporter)

This is a server (aka Prometheus-compatible exporter) which checks the configured hosts against various DNSBL (DNS Block Lists), sometimes referred to as RBLs.

Should you accept this mission, your task is to scrape `/metrics` using Prometheus to create graphs, alerts, and so on.


> [!IMPORTANT]
> This is (still) pretty early software. But I happily accept all kinds of feedback - bug reports, PRs, code, docs, ... :)

> [!TIP]
> Documentation is available in [docs](./docs/src/).

## License / Author

This code is Apache 2.0 licensed.

For questions, comments or anything else, [please get in touch](https://www.luzilla-capital.com).

## Releasing new versions

(This is for myself, since I tend to forget things.)

 1. `git tag -a x.y.z`
 1. `git push --tags`
 1. GitHub Actions/GoReleaser will build a pretty release
