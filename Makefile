.PHONY: build
build:
	goreleaser build --snapshot --single-target --clean

.PHONY: run-dev
run-dev:
	go run dnsbl_exporter.go --log.debug

.PHONY: snapshot
snapshot:
	goreleaser build --snapshot --clean

.PHONY: test
test:
	act "pull_request" -j test
