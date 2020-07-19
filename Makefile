# Image URL to use all building/pushing image targets
IMG ?= evalsocket/oldmonk:latest

VERSION ?= $(shell cat ./VERSION)
PKG_PATH ?= github.com/remotegarage/oldmonk
BUILDTIME := $(shell date -u +%Y%m%d.%H%M%S)
LDFLAGS ?= -ldflags '-X ${PKG_PATH}/pkg/version.version=${VERSION} -X ${PKG_PATH}/pkg/version.buildtime=${BUILDTIME}'

all: test manager

# Install all the build and lint dependencies
setup:
	go mod download
	go generate -v ./...
.PHONY: setup

# Run tests
test:
	LC_ALL=C go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

# Build manager binary
build: mod fmt 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/_output/bin/oldmonk $(PKG_PATH)/cmd/manager

# Build manager binary
ci: mod generate fmt  build


docker-manager:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/manager $(PKG_PATH)/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: mod  manifests generate fmt 
	go run ./cmd/manager/main.go


# Install CRDs into a cluster
install:
	kustomize build ./deploy | kubectl apply -f -
	
# Generate manifests e.g. CRD, RBAC etc.
manifests:
	operator-sdk generate k8s
	operator-sdk generate crds
	go run github.com/ahmetb/gen-crd-api-reference-docs -api-dir ./pkg/apis -config docs/api/config.json -template-dir docs/api/template -out-file docs/api/index.html


# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/... ./x/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

# Get modendencies
mod:
	go mod tidy

# Build the docker image
docker-build: mod  manifests generate fmt  test build
	operator-sdk build ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i "" 's|REPLACE_IMAGE|${IMG}|g' deploy/operator.yaml

# Push the docker image
docker-push:
	docker push ${IMG}
