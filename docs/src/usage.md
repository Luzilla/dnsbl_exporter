# Usage

Learn how to configure the `dnsbl-exporter`.

## Configuration

See `rbls.ini` and `targets.ini` files in the repository for examples.

The files follow the Nagios format as this exporter is meant to be a drop-in replacement so you can factor out Nagios, one (simple) step at a time. 😊

### Checking both IP and domain based RBLs

By default `--config.rbls` is checked as an IP based blocklist (targets are
resolved to IPs first), or, when `--config.domain-based` is set, as a domain
based blocklist (targets are used as-is).

To check both from the same deployment, set `--config.rbls-domain` to a file
containing domain based RBLs (e.g. `dbl.spamhaus.org`); `--config.rbls` then
stays IP based. Every target from `--config.targets` is checked against
**both** lists, unless you also set `--config.targets-domain` to give the
domain based checklist its own, separate target list. Without it, an IP
literal in `--config.targets` gets a (harmless but pointless) domain-style
lookup, and every domain based target also gets resolved and checked
against the IP based RBLs — so set `--config.targets-domain` whenever the
two checklists are meant to cover different targets.

Otherwise:

```sh
$ dnsbl-exporter -h
...
--config.dns-resolver value  IP address of the resolver to use. (default: "127.0.0.1:53")
--config.rbls value          Configuration file which contains RBLs (default: "./rbls.ini")
--config.rbls-domain value   Configuration file which contains domain based RBLs, checked in addition to --config.rbls
--config.targets value       Configuration file which contains the targets to check. (default: "./targets.ini")
--config.targets-domain value  Configuration file which contains the targets to check against --config.rbls-domain. Defaults to --config.targets when unset.
--config.domain-based        RBLS are domain instead of IP based blocklists (default: false)
--web.listen-address value   Address to listen on for web interface and telemetry. (default: ":9211")
--web.telemetry-path value   Path under which to expose metrics. (default: "/metrics")
--log.debug                  Enable more output in the logs, otherwise INFO.
--log.output value           Destination of our logs: stdout, stderr (default: "stdout")
--help, -h                   show help
--version, -V                Print the version information.
```

### System resolver

The `dnsbl-exporter` can use your system resolver from `/etc/resolv.conf` automatically.


> **Please note:**
> The `dnsbl-exporter` needs read permissions to `/etc/resolv.conf` file for this feature to work.

Configure resolver as an argument:

```sh
--config.dns-resolver=system
```

Configure resolver as an environment variable:

```dotenv
DNSBL_EXP_RESOLVER=system
```

<div class="warning">
Please see <strong>DNS</strong> for further details.

Not every resolver is compatible with most RBLs.
</div>

## Metrics returned by exporter

The individual configured servers and their status are represented by a **gauge**:

```sh
luzilla_rbls_ips_blacklisted{hostname="mail.gmx.net",ip="212.227.17.168",rbl="ix.dnsbl.manitu.net"} 0
```

This represents the server's hostname and the DNSBL in question:

 - `0` (zero) for unlisted
 - `1` (one) for listed

Requests to the DNSBL happen in real-time and are not cached. Take this into account and use accordingly.

If the exporter is configured for DNS based blocklists, the ip label represents the return code of the blocklist.

<div class="warning">
You are listed!

If you happen to be listed — inspect the exporter's logs as they will contain a reason.
</div>