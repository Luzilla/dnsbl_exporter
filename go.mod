module github.com/luzilla/dnsbl_exporter

go 1.12

require github.com/prometheus/client_golang v1.9.0

require (
	github.com/Luzilla/godnsbl v1.0.0
	github.com/miekg/dns v1.1.37
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/urfave/cli v1.22.5
	gopkg.in/ini.v1 v1.62.0
)
