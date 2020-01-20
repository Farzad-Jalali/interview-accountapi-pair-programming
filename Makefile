.DEFAULT_GOAL := default
GO111MODULE := off
platform := $(shell uname)
swagger_codegen_version := "v0.18.0"
secscan_image := "288840537196.dkr.ecr.eu-west-1.amazonaws.com/tech.form3/secscan-go:latest"

ifeq (${platform},Darwin)
swagger_binary := "swagger_darwin_amd64"
else
swagger_binary := "swagger_linux_amd64"
endif

GOFMT_FILES?=$$(find ./ -name '*.go' | grep -v vendor)

default: install-deps build test

build: goimportscheck errcheck vet
	@find ./cmd/* -maxdepth 1 -type d -exec go install "{}" \;

install-deps: install-goimports

install-swagger:
	@curl -o /usr/local/bin/swagger -L'#' https://github.com/go-swagger/go-swagger/releases/download/${swagger_codegen_version}/${swagger_binary} && chmod +x /usr/local/bin/swagger

install-lint:
	@go get -u golang.org/x/lint/golint

install-goimports:
	@if [ ! -f ./goimports ]; then \
		go get golang.org/x/tools/cmd/goimports; \
	fi

test: goimportscheck
	@echo "executing tests..."
	@go test -v -count 1 github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

goimports:
	goimports -w $(GOFMT_FILES)

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

docker-package: build
	@find ./cmd/* -maxdepth 1 -type d -exec sh -c '"$(CURDIR)/scripts/docker-package.sh" {}' \;

docker-publish: docker-package
	@find ./cmd/* -maxdepth 1 -type d -exec sh -c '"$(CURDIR)/scripts/docker-publish.sh" {}' \;

publish: docker-publish

lint:
	@echo "go lint ."
	@golint $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Lint found errors in the source code. Please check the reported errors"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

swagger:
	@exec sh -c '"$(CURDIR)/scripts/swagger.sh" {}'

secscan:
	@exec sh -c 'eval $$(aws ecr get-login --region eu-west-1 --no-include-email) && docker pull ${secscan_image}'
	@exec docker run --rm -v $(CURDIR):/code -e TRAVIS -e REPO=form3tech/interview-accountapi -e SNYK_TOKEN=${SNYK_TOKEN} ${secscan_image}

run:
	@docker-compose run --rm wait_for
	@docker-compose up accountapi

docs:
	@docker run -v $$PWD/:/docs pandoc/latex -f markdown /docs/CANDIDATE_INSTRUCTIONS.md -o /docs/build/output/instructions.pdf

.PHONY: build test testacc vet goimports goimportscheck errcheck docker-package lint docker-publish vendor-status test-compile swagger
