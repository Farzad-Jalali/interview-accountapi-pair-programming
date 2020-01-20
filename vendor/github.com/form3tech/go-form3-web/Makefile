.DEFAULT_GOAL := default

platform := $(shell uname)

GOFMT_FILES?=$$(find ./ -name '*.go' | grep -v vendor)

default: build test

build: goimportscheck vet
	go install ./httpclient

install-goimports:
	@go get golang.org/x/tools/cmd/goimports

test: goimportscheck
	@echo "executing tests..."
	@go test github.com/form3tech/go-form3-web/httpclient

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

swagger-generate-paymentapi-client:
	@test "${GITHUB_TOKEN}" || (echo 'you must define GITHUB_TOKEN env variable' && exit 1)
	@mkdir -p ./httpclient/support/swagger-client/paymentapi
	@bash ./scripts/gh-download-release.sh "$(GITHUB_TOKEN)" form3tech paymentapi swagger.yaml latest
	@swagger generate client -f ./httpclient/support/swagger-client/paymentapi/swagger.yaml -t ./httpclient/support/swagger-client/paymentapi

goimports:
	goimports -w $(GOFMT_FILES)

goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"


.PHONY: build test vet goimports goimportscheck errcheck package
