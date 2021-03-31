.PHONY: build
build:
	goreleaser --snapshot --rm-dist

.PHONY: test
test:
	docker run \
		-it \
		--rm \
		-v $(CURDIR):/src/github.com/Luzilla/dnsbl_exporter \
		-w /src/github.com/Luzilla/dnsbl_exporter \
		golang:1.16 \
		sh -c "go mod download && go test ./..."
