release_name:=dnsbl-exporter-dev
namespace:=test

.PHONY: build
build:
	goreleaser build --snapshot --single-target --clean

.PHONY: clean
clean:
	rm -rf dist/
	rm ./chart/*.tgz

.PHONY: run-dev
run-dev:
	go run dnsbl_exporter.go \
		--log.debug \
		--config.dns-resolver 0.0.0.0:15353

.PHONY: run-dev-domain
run-dev-domain:
	go run dnsbl_exporter.go \
		--log.debug \
		--config.rbls ./rbls-domain.ini \
		--config.domain-based \
		--config.dns-resolver 0.0.0.0:15353

.PHONY: snapshot
snapshot:
	goreleaser build --snapshot --clean

.PHONY: test
test:
	go test -v --race ./...

.PHOMY: test-deploy-helm
test-deploy-helm:
	helm upgrade --install --namespace $(namespace) \
		-f ./chart/values.yaml -f ./chart/values.dev.yml -f ./chart/values.domain-based.yaml \
		$(release_name) ./chart

.PHONY: test-undeploy-helm
test-undeploy-helm:
	helm uninstall -n $(namespace) $(release_name)

.PHONY: build-unbound
build-unbound:
	docker build \
		-t ghcr.io/luzilla/unbound:dev \
		.docker/unbound/rootfs

# Dalec package builds (requires goreleaser snapshot first)
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")
# GoReleaser arch directory suffix: amd64_v1, arm64_v8.0
ARCH_DIR     ?= amd64_v1
BINARY_PATH  ?= dist/dnsbl_exporter_linux_$(ARCH_DIR)/dnsbl-exporter

DALEC_TARGETS ?= mariner2/rpm jammy/deb noble/deb bookworm/deb azlinux3/testing/sysext
DALEC_TARGET  ?= mariner2/rpm

dalec_out = $(lastword $(subst /, ,$(1)))

.PHONY: dalec
dalec:
	$(info Building $(DALEC_TARGET) for $(VERSION) ($(BINARY_PATH)))
	docker build -f .dalec/dalec.yml \
		--platform=linux/amd64 \
		--target=$(DALEC_TARGET) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BINARY_PATH=$(BINARY_PATH) \
		-o type=local,dest=./out/$(call dalec_out,$(DALEC_TARGET)) .

.PHONY: dalec-all
dalec-all: snapshot
	@for t in $(DALEC_TARGETS); do \
		$(MAKE) dalec DALEC_TARGET=$$t || exit 1; \
	done

.PHONY: dalec-debug
dalec-debug:
	@docker run --rm -v "$(CURDIR)/out/sysext:/sysext" alpine:latest sh -c '\
		apk add -q erofs-utils && \
		for f in /sysext/*.raw; do \
			echo "=== $$f ==="; \
			dump.erofs --ls "$$f"; \
			echo "--- extension-release ---"; \
			fsck.erofs --extract=/tmp/x "$$f" >/dev/null && \
			cat /tmp/x/usr/lib/extension-release.d/extension-release.* ; \
			rm -rf /tmp/x; \
		done'
