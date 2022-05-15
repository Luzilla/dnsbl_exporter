GO_VERSION:=1.16

.PHONY: build
build:
	goreleaser build --snapshot --single-target --rm-dist

.PHONY: snapshot
snapshot:
	goreleaser build --snapshot --rm-dist

.PHONY: test
test:
	docker run \
		-it \
		--rm \
		-v $(CURDIR):/src/github.com/Luzilla/dnsbl_exporter \
		-w /src/github.com/Luzilla/dnsbl_exporter \
		golang:$(GO_VERSION) \
		sh -c "go mod download && go test ./..."
