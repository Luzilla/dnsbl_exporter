# DNSBL Exporter

The helm chart installs DNSBL Exporter into Kubernetes cluster. 

## Installation

Helm chart depends on `bitnami-common` helm chart from [Bitnami](https://github.com/bitnami/charts).
You need to install the repo before you can make any deployment:
```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Installation/Update looks like:

```shell
helm repo update bitnami
helm dependencies build

# Create namespace (only once)
kubectl create namespace dnsbl

# Install/Update
helm upgrade dnsbl ./ \
    --install \
    --create-namespace \
    --namespace dnsbl ${OPTS} \
    --values ./values.dnsbl.yml
```

Here we rely on `values.dnsbl.yml` which overrides the defaults in chart's `values.yaml`. As usual the defaults work well,
but if you would like to customise the setup, the file may look like this:
```yaml
---
image:
  registry: docker.ti.local
  repository: dnsbl-exporter
  tag: 0.6.0
  pullSecrets:
    - rs-docker-registry

pdb:
  create: true
  minAvailable: 1
  maxUnavailable: ""

service:
  type: ClusterIP

resources:
  limits:
    cpu: 250m
    memory: 256Mi
  requests:
    cpu: 50m
    memory: 16Mi

configuration:
  rbls:
    # 0SPAM
    - bl.0spam.org
    - rbl.0spam.org
    # NIXSPAM
    - ix.dnsbl.manitu.net
    # Spamhaus Zen
    - zen.spamhaus.org
    # Anonmails DNSBL
    - spam.dnsbl.anonmails.de
    # UCEPROTECT-Level 1
    - dnsbl-1.uceprotect.net
    # UCEPROTECT-Level 2
    - dnsbl-2.uceprotect.net
    # UCEPROTECT-Level 2
    - dnsbl-3.uceprotect.net
    # Backscatterer
    - ips.backscatterer.org
    # Barracuda Reputation Block List
    - b.barracudacentral.org
    # Blocklist.de
    - bl.blocklist.de
    # CALIVENT
    - dnsbl.calivent.com.pe
    # CYMRU BOGONS
    - bogons.cymru.com
    # DNS Servicios
    - rbl.dns-servicios.com
    # DRMX
    - bl.drmx.org
    # DRONEBL
    - dnsbl.dronebl.org
    # FABEL SOURCES
    - spamsources.fabel.dk
    # HIL HABEAS
    - hil.habeas.com
    # HIL2 HABEAS
    - hil2.habeas.com
    # Hostkarma
    - hostkarma.junkemailfilter.com
    # IBM DNS Blacklist
    - dnsbl.cobion.com
    # ICM FORBIDDEN
    - forbidden.icm.edu.pl
    # IMP WORM
    - dnsrbl.swinog.ch
    # IMP SPAM
    - spamrbl.swinog.ch
    - uribl.swinog.ch
    # Spamhaus ZEN
    - zen.spamhaus.org
    # Spamhaus DBL - should not be used with IPs
    #- dbl.spamhaus.org
    - xbl.spamhaus.org
    # SPFBL DNSBL
    - dnsbl.spfbl.net
    # Sender Score Reputation Network
    - bl.score.senderscore.com
    # SORBS BLOCK
    - block.dnsbl.sorbs.net
    # SORBS DUHL
    - dul.dnsbl.sorbs.net
    # SORBS HTTP
    - http.dnsbl.sorbs.net
    # SORBS MISC
    - misc.dnsbl.sorbs.net
    # SORBS NEW
    - new.spam.dnsbl.sorbs.net
    # SORBS SMTP
    - smtp.dnsbl.sorbs.net
    # SORBS SOCKS
    - socks.dnsbl.sorbs.net
    # SORBS SPAM
    - spam.dnsbl.sorbs.net
    # SORBS WEB
    - web.dnsbl.sorbs.net
    # SORBS ZOMBIE
    - zombie.dnsbl.sorbs.net
    # RATS Dyna
    - dyna.spamrats.com
    # RATS NoPtr
    - noptr.spamrats.com
    # RATS Spam
    - spam.spamrats.com
    # SEM BACKSCATTER
    - backscatter.spameatingmonkey.net
    # SEM BLACK
    - bl.spameatingmonkey.net
    # MSRBL Phishing
    - phishing.rbl.msrbl.net
    # MSRBL Spam
    - spam.rbl.msrbl.net
    # NETHERRELAYS
    - relays.nether.net
    # NETHERUNSURE
    - unsure.nether.net
    # NIXSPAM
    - ix.dnsbl.manitu.net
    # Nordspam BL
    - bl.nordspam.com
    # NoSolicitado
    - bl.nosolicitado.org
    # ORVEDB
    - orvedb.aupads.org
    # PSBL
    - psbl.surriel.com
    # RBL JP
    - virus.rbl.jp
    # RSBL
    - rsbl.aupads.org
    # s5h.net
    - all.s5h.net
    # SCHULTE
    - rbl.schulte.org
    # SERVICESNET
    - korea.services.net
    # SPAMCOP
    - bl.spamcop.net
    # Suomispam Reputation
    - bl.suomispam.net
    # SWINOG
    - dnsrbl.swinog.ch
    # TRIUMF
    - rbl2.triumf.ca
    # TRUNCATE
    - truncate.gbudb.net
    # Woodys SMTP Blacklist
    - blacklist.woody.ch
    # WPBL
    - db.wpbl.info
    # ZapBL
    - dnsbl.zapbl.net
    # INTERSERVER
    - rbl.interserver.net
    # JIPPG
    - dialup.blacklist.jippg.org
    # KEMPTBL
    - dnsbl.kempt.net
    # KISA
    - spamlist.or.kr
    # Konstant
    - bl.konstant.no
    # LASHBACK
    - ubl.lashback.com
    # LNSGBLOCK
    - spamguard.leadmon.net
    # MADAVI
    - dnsbl.madavi.de
    # MAILSPIKE BL
    - bl.mailspike.net
    # MAILSPIKE Z
    - z.mailspike.net
  targets:
    # Main Server
    - hq.example.com
    # Backup MX
    - mxb.example.com
    # Sendgrid dedicated IP
    - 10.1.1.1
  resolver: "10.2.3.4:53"

extraDeploy:
  - apiVersion: monitoring.coreos.com/v1
    kind: PrometheusRule
    metadata:
      name: dnsbl-rules
    spec:
      groups:
        - name: dnsbl
          rules:
            - alert: RblsIpsBlacklisted
              expr: max_over_time(luzilla_rbls_ips_blacklisted[60m]) > 0
              for: 60m
              labels:
                severity: high
              annotations:
                description: '{{ "{{" }} $labels.hostname {{ "}}" }} ({{ "{{" }} $labels.ip {{ "}}" }}) has been blacklisted in {{ "{{" }} $labels.rbl {{ "}}" }} for more than 60 minutes.'
                summary: 'Endpoint {{ "{{" }} $labels.hostname {{ "}}" }} is blacklisted'
```
