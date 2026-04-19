# Container

Docker/OCI images are available in the [container registry](https://github.com/orgs/Luzilla/packages?repo_name=dnsbl_exporter):

```sh
$ docker pull ghcr.io/luzilla/dnsbl_exporter:vX.Y.Z
...
```

The image works with podman and docker as the exporter runs rootless (user `nobody`).

> [!CAUTION]
> The `latest` tag is not provided. Please pick an explicit version.

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
    -v ./conf:/etc/dnsbl-exporter \
    ghcr.io/luzilla/dnsbl_exporter:vA.B.C
```

```Dockerfile
FROM ghcr.io/luzilla/dnsbl_exporter:vA.B.C

ADD my-target.ini /target.ini
ADD my-rbls.ini /rbls.ini
```

## Helm

Additionally, a helm chart is provided to run the `dnsbl-exporter` on Kubernetes.

To get started quickly, an unbound container is included as a sidecar in the pod alongside the exporter.

This unbound acts as a local DNS server to send queries to.

> [!TIP]
> To disable the included, use `unbound.enabled=false` and provide your own resolver (via `config.resolver: an.ip.address:port`).

To configure the chart, copy [`chart/values.yaml`](https://github.com/Luzilla/dnsbl_exporter/blob/main/chart/values.yaml) to `values.local.yaml`.

Another useful option to add our chart as dependency to your own chart:

```yaml
dependencies:
  - name: dnsbl-exporter
    repository: oci://ghcr.io/luzilla/charts
    version: 0.1.0
```

> [!CAUTION]
> Please adjust the chart version.

The sources for the helm chart are in [chart](https://github.com/Luzilla/dnsbl_exporter/tree/main/chart), to install it, you can inspect the `Chart.yaml` for the version, check the [helm chart repository](https://github.com/orgs/Luzilla/packages/container/package/charts%2Fdnsbl-exporter) or check out [artifact hub](https://artifacthub.io/packages/helm/luzilla/dnsbl-exporter).

The following command creates a `dnsbl-exporter` release which is installed into a namespace called `my-namespace`:

```sh
helm upgrade --install \
    --namespace my-namespace \
    -f ./chart/values.yaml \
    -f ./values.local.yaml \
    dnsbl-exporter oci://ghcr.io/luzilla/charts/dnsbl-exporter --version 0.1.0
```