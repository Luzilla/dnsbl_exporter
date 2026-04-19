# Standalone

_Releases are available for Mac, Linux and FreeBSD._

 1. Go to [release](https://github.com/Luzilla/dnsbl_exporter/releases) and download a binary for your system/platform.
 1. Get `rbls.ini` and put it next to the binary.
 1. Get `targets.ini`, and customize. Or use the defaults.
 1. `./dnsbl-exporter`

> [!TIP]
> Example config files are provided on [GitHub](https://github.com/Luzilla/dnsbl-exporter).

Go to `http://127.0.0.1:9211/` in your browser.

## Managing the exporter

> [!TIP]
> Check out [Packages](./packages.md).

As option you can configure exporter to run as `systemd` service. Consider using a package instead.