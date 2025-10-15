.PHONY: build test test-unit test-integration test-conformance docker-build docker-push deploy clean lint vet fmt verify

GO_VERSION := 1.21
IMAGE_NAME := cert-manager-webhook-joker
IMAGE_TAG := latest
REGISTRY := ghcr.io/nauski
VERSION := v1.0.1

build:
	go build -o webhook ./cmd/webhook

test:
	go test -v ./...

test-unit:
	go test -v ./internal/...

test-integration:
	go test -v ./test/integration/...

test-conformance:
	@if [ -z "$(TEST_ZONE_NAME)" ]; then \
		echo "TEST_ZONE_NAME must be set for conformance tests"; \
		echo "Example: export TEST_ZONE_NAME=test.example.com"; \
		exit 1; \
	fi
	@if [ -z "$(TEST_SECRET_NAME)" ]; then \
		export TEST_SECRET_NAME=joker-credentials; \
	fi
	go test -v ./test/conformance/...

lint:
	golangci-lint run

vet:
	go vet ./...

fmt:
	go fmt ./...

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):latest

docker-push: docker-build
	docker push $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)/$(IMAGE_NAME):latest

deploy:
	kubectl apply -f deploy/

deploy-secret:
	kubectl apply -f deploy/secret.yaml

undeploy:
	kubectl delete -f deploy/ --ignore-not-found=true

clean:
	rm -f webhook
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) || true

deps:
	go mod download
	go mod tidy

verify: fmt vet lint test

.DEFAULT_GOAL := build
