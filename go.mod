module github.com/luzilla/dnsbl_exporter

go 1.12

require github.com/prometheus/client_golang v1.2.1

require (
	github.com/Luzilla/godnsbl v1.0.0
	github.com/miekg/dns v1.1.22
	github.com/prometheus/common v0.7.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/urfave/cli v1.22.1
	gopkg.in/ini.v1 v1.49.0
)
