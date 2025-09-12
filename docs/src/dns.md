# DNS

In order to use the `dnsbl-exporter`, you need a DNS resolver to query the RBLs directly.

Using public resolvers such as Google, Cloudflare, OpenDNS, Quad9 etc. will likely not work as often a `NXDOMAIN` response is treated different by these providers.

The recommendation is to setup a private resolver, such as [Unbound](https://github.com/NLnetLabs/unbound) without forwarding, which means you will use root NS.

> **Please note:** If you have local copy (mirror) of RBL zone synced over rsync or other channels you can configure local [`rbldnsd`](https://rbldnsd.io/) to serve this zone and configure Unbound query this zone from `rbldnsd`.

To install unbound on a Mac, follow these steps:

```sh
$ brew install unbound
...
$ sudo unbound -d -vvvv
```

> And leave the Terminal open — there will be ample queries and data for you to see and learn from.

An alternative to Homebrew is to use Docker; an example image is provided in this repository, it
contains a working configuration — ymmv.

```sh
docker run -p 53:5353/udp ghcr.io/luzilla/unbound:v0.7.0-rc3
```

Verify Unbound works:

```sh
 $ dig +short @127.0.0.1 spamhaus.org
192.42.118.104
```