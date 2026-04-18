# packages

> [!TIP]
> New in version `v0.13.0`.

Each (GitHub) release also contains:

 - rpm
 - deb (jammy, noble, bookworm)
 - sysext

The package will setup default configuration in `/etc/dnsbl-exporter` and includes a systemd unit to manage the exporter as well.

```bash
$ systemctl enable dnsbl-exporter
$ systemctl start dnsbl-exporter
$ systemctl status dnsbl-exporter
```