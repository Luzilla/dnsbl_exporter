GO_VERSION:=0.16

.PHONY: build
build:
	goreleaser r --snapshot --rm-dist

.PHONY: test
test:
	docker run \
		-it \
		--rm \
		-v $(CURDIR):/src/github.com/Luzilla/dnsbl_exporter \
		-w /src/github.com/Luzilla/dnsbl_exporter \
		golang:$(GO_VERSION) \
		sh -c "go mod download && go test ./..."
