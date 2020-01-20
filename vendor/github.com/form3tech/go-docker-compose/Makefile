.DEFAULT_GOAL := default

platform := $(shell uname)

GOFMT_FILES?=$$(find ./ -name '*.go' | grep -v vendor)

default: build test

build: goimportscheck vet
	go install ./dockercompose

install-goimports:
	@go get golang.org/x/tools/cmd/goimports

test: goimportscheck
	@echo "executing tests..."
	@go test -count 1 -v github.com/form3tech/go-docker-compose/dockercompose

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

goimports:
goimports:
	goimports -w $(GOFMT_FILES)

goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"


.PHONY: build test vet goimports goimportscheck errcheck package
