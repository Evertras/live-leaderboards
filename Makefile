.PHONY: default
default: .git/hooks/pre-commit generated bin/terraform
	$(MAKE) -C site generated
	@echo Ready!

.PHONY: build
build: build-api build-site

.PHONY: build-api
build-api: bin/leaderboard-api-lambda

.PHONY: build-site
build-site:
	$(MAKE) -C site build

.PHONY: validate-schema
validate-schema: node_modules
	npx @redocly/cli lint ./specs/openapi.yaml

.PHONY: terraform-apply
terraform-apply: bin/terraform ./bin/leaderboard-api-lambda
	cd terraform && ../bin/terraform apply

.PHONY: clean
clean:
	rm -rf node_modules
	$(MAKE) -C site clean

################################################################################
# Generated stuff
.PHONY: generated
generated: ./pkg/api/api.go diagrams node_modules
	$(MAKE) -C site generated

GO_FILES=$(shell find . -iname *.go)
bin/leaderboard-api-lambda: $(GO_FILES) go.mod go.sum
	@mkdir -p bin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/leaderboard-api-lambda -ldflags="-s -w" ./cmd/lambda

bin/leaderboard-api-lambda.zip: bin/leaderboard-api-lambda
	cd bin && zip leaderboard-api-lambda.zip leaderboard-api-lambda

./pkg/api/api.go: specs/openapi.yaml bin/oapi-codegen
	@mkdir -p pkg/api
	./bin/oapi-codegen -package api specs/openapi.yaml > pkg/api/api.go

DIAGRAM_PUML=$(shell find docs/diagrams -iname '*.puml')
DIAGRAM_SVG=$(DIAGRAM_PUML:.puml=.svg)

.PHONY: diagrams
diagrams: $(DIAGRAM_SVG)

# TODO: Better generation that doesn't require global plantuml install
%.svg: %.puml
	plantuml -svg $<

################################################################################
# Local dependencies
node_modules: package.json package-lock.json
	npm install
	@touch node_modules

################################################################################
# Testing
test: ./pkg/api/api.go
	go test -race ./pkg/...

.PHONY: test-integration
test-integration: ./pkg/api/api.go
	go test -v -race ./tests

.PHONY: docker-compose-up
docker-compose-up: ./pkg/api/api.go
	cd tests && docker-compose up --build

################################################################################
# Formatting
#
# Commands to format stuff to standards

# Format everything
.PHONY: fmt
fmt: fmt-terraform fmt-prettier

.PHONY: fmt-prettier
fmt-prettier: node_modules
	npx prettier . --write

.PHONY: fmt-terraform
fmt-terraform: bin/terraform
	terraform fmt -recursive ./terraform

# Run fmt on pre-commit
.git/hooks/pre-commit: .evertras/pre-commit.sh
	cp .evertras/pre-commit.sh .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

################################################################################
# Swagger for API specs
#
# Runs Swagger to view our API spec
.PHONY: swagger
swagger:
	@echo "Hosting at http://localhost:8080"
	docker run -p 8080:8080 \
		-e SWAGGER_JSON=/data/openapi.yaml \
		-v ./specs/:/data/ \
		swaggerapi/swagger-ui:v5.4.2

################################################################################
# Local tooling
#
# This section contains tools to download to the local ./bin directory for easy
# local use.  The .envrc file makes adding the local ./bin directory to our path
# simple, so we can use tools here without having to install them globally as if
# they actually were global.

# For now we only support Linux 64 bit and MacOS for simplicity
ifeq ($(shell uname), Darwin)
OS_URL := darwin
else
OS_URL := linux
endif

TERRAFORM_VERSION=1.5.5
bin/terraform:
	@mkdir -p bin
	curl -Lo bin/terraform.zip https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_$(OS_URL)_amd64.zip
	cd bin && unzip terraform.zip
	rm bin/terraform.zip

bin/oapi-codegen:
	GOBIN=$(shell pwd)/bin go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.4
