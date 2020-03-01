# Image URL to use all building/pushing image targets
IMG ?= gcr.io/multi-k8s-259909/oldmonk:latest

VERSION ?= $(shell cat ./VERSION)
PKG_PATH ?= github.com/evalsocket/oldmonk
BUILDTIME := $(shell date -u +%Y%m%d.%H%M%S)
LDFLAGS ?= -ldflags '-X ${PKG_PATH}/pkg/version.version=${VERSION} -X ${PKG_PATH}/pkg/version.buildtime=${BUILDTIME}'

all: test manager

# Run tests
test: mod generate fmt vet manifests
	go test $(LDFLAGS) ./pkg/... ./cmd/... -coverprofile cover.out

# Build manager binary
manager: mod manifests generate fmt vet
	go build $(LDFLAGS) -o bin/manager $(PKG_PATH)/cmd/manager

# Build manager binary
ci: mod generate fmt vet
		go build $(LDFLAGS) -o bin/manager $(PKG_PATH)/cmd/manager


docker-manager:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/manager $(PKG_PATH)/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: mod  manifests generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install:
	kubectl apply -f deploy/crds
	kubectl apply -f deploy/service_account.yaml
	kubectl apply -f deploy/role.yaml
	kubectl apply -f deploy/role_binding.yaml
	kubectl apply -f deploy/operator.yaml


# Generate manifests e.g. CRD, RBAC etc.
manifests:
	operator-sdk generate k8s
	operator-sdk generate crds
	go run github.com/ahmetb/gen-crd-api-reference-docs -api-dir ./pkg/apis -config docs/api/config.json -template-dir docs/api/template -out-file docs/api/index.html


# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/... ./x/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

# Get modendencies
mod:
	go mod tidy

# Build the docker image
docker-build: mod  manifests generate fmt vet test
	operator-sdk build ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i "" 's|REPLACE_IMAGE|${IMG}|g' deploy/operator.yaml

# Push the docker image
docker-push:
	docker push ${IMG}
