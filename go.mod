module github.com/luzilla/dnsbl_exporter

go 1.12

require github.com/prometheus/client_golang v1.2.1

require (
	github.com/Luzilla/godnsbl v1.0.0
	github.com/miekg/dns v1.1.26
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/urfave/cli v1.22.1
	gopkg.in/ini.v1 v1.49.0
)
