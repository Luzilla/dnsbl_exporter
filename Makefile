release_name:=dnsbl-exporter-dev
namespace:=test

.PHONY: build
build:
	goreleaser build --snapshot --single-target --clean

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
	act "pull_request"

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
