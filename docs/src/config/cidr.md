# CIDR targets

`targets.ini` accepts individual IP addresses, hostnames, *and* IPv4 CIDR
ranges. The exporter expands each CIDR into its constituent host
addresses at scrape time and queries every configured RBL for each one.

```ini
[targets]
server=8.8.8.8         ; single IP
server=1.1.1.1/32      ; CIDR — single host
server=192.0.2.0/24    ; CIDR — 254 hosts (network and broadcast excluded)
server=mail.example.org ; hostname (resolved to A records, then queried)
```

## How expansion works

For a given prefix the exporter emits one query per RBL per host:

| prefix | hosts produced | notes |
|--------|---------------:|-------|
| `/32`  | 1   | the address itself |
| `/31`  | 2   | both addresses (RFC 3021) |
| `/30`  | 2   | network and broadcast excluded |
| `/29`  | 6   | network and broadcast excluded |
| `/24`  | 254 | typical LAN scan |
| `/22`  | 1022 |  |
| `/16`  | 65 534 | |

Cost per scrape is roughly **`expanded_hosts × configured_RBLs`** DNS
queries. Most public RBLs apply per-resolver query rate limits; large
expansions will trip them quickly and your scrapes will start returning
errors instead of useful data.

## Large-prefix warning

> [!TIP]
> New in `v0.13.0`.

When a single CIDR expands to **1024 or more** hosts (i.e. `/21` or larger
on IPv4), the exporter emits a `WARN` log record on every scrape:

```
level=WARN msg="CIDR expands to a large number of hosts; scrape load and resolver rate-limiting risk" cidr=10.0.0.0/16 hosts=65534
```

The warning is informational — the expansion still happens and the scrape
still runs. The exporter does not refuse, truncate, or cap large
expansions; that decision is left to the operator.

> [!WARNING]
> The warn fires **once per scrape**, not once at startup. With Prometheus'
> default 15s scrape interval that is four warn records per minute, per
> oversized CIDR. This is intentional — silence it by shrinking the prefix
> in `targets.ini` rather than filtering the log.

## Practical recommendations

- For monitoring a handful of mail relays or outbound IPs, list them
  individually or use small CIDRs (`/29`, `/28`). The warning leaves you
  alone, the scrapes stay fast, and the RBLs don't notice you.
- `/24` is a reasonable upper bound for "monitor everything in this
  subnet." Below that you start consuming a meaningful fraction of most
  RBLs' free-tier rate limits within a single scrape.
- Anything `/22` or larger should be a deliberate decision — and probably
  paired with a private mirror of the RBL zones (rsync / rbldnsd) so you
  aren't hitting the public infrastructure.
- Avoid `/16` and larger entirely.

> [!NOTE]
> CIDR expansion is **IPv4 only** today. IPv6 prefixes in `targets.ini`
> will not behave correctly — track the upstream issue if this matters
> for your deployment.
