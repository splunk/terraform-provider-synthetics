TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=splunk
NAME=synthetics
BINARY=terraform-provider-${NAME}
VERSION=3.0.0-dev

.PHONY: default tools fmt fmtcheck lint vet govulncheck docs vendor-check build install test test-race testacc

default: install

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./${NAME}/*.go

fmtcheck:
	@echo "==> Checking source code is gofmt'd..."
	@out=$$(gofmt -s -l ./${NAME}/*.go ./main.go); if [ -n "$$out" ]; then \
		echo "The following files are not gofmt'd:"; echo "$$out"; exit 1; \
	fi

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

vet:
	@echo "go vet ."
	@go vet -mod=vendor $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

govulncheck:
	@echo "==> Checking for known vulnerabilities..."
	govulncheck ./...

docs:
	@echo "==> Generating provider documentation..."
	go generate ./...

vendor-check:
	@echo "==> Checking vendor/ is in sync with go.mod..."
	go mod tidy
	go mod vendor
	git diff --exit-code -- go.mod go.sum vendor/

build:
	go build -mod=vendor -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -mod=vendor $(TEST) $(TESTARGS) -timeout=30s -parallel=4

test-race:
	go test -mod=vendor -race $(TEST) $(TESTARGS) -timeout=60s -parallel=4

testacc: SHELL:=/bin/bash
testacc:
	set -o pipefail; TF_ACC=1 go test -json $(TEST) $(TESTARGS) -timeout 30m \
		| sed '/X-Sf-Token/d' \
		| tee testacc.jsonl \
		| jq -j -r 'if .Action == "output" then .Output else empty end'