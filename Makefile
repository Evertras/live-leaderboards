.PHONY: default
default: .git/hooks/pre-commit node_modules
	@# Do nothing

.PHONY: test-integration
test-integration:
	go test -v -race ./tests

DIAGRAM_PUML=$(shell find docs/diagrams -iname '*.puml')
DIAGRAM_SVG=$(DIAGRAM_PUML:.puml=.svg)

.PHONY: diagrams
diagrams: $(DIAGRAM_SVG)

# TODO: Better generation that doesn't require global plantuml install
%.svg: %.puml
	plantuml -svg $<

# Format everything
.PHONY: fmt
fmt: bin/terraform
	npx prettier . --write
	terraform fmt -recursive ./terraform

# Prettier installs node_modules
node_modules: package.json package-lock.json
	npm install

# Run prettier on pre-commit
.git/hooks/pre-commit:
	cp .evertras/pre-commit.sh .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

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
