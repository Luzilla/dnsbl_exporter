.PHONY: build
build:
	goreleaser --snapshot --rm-dist
