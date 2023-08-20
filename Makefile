.PHONY: default
default: .git/hooks/pre-commit node_modules
	@# Do nothing

DIAGRAM_PUML=$(shell find docs/diagrams -iname '*.puml')
DIAGRAM_SVG=$(DIAGRAM_PUML:.puml=.svg)

.PHONY: diagrams
diagrams: $(DIAGRAM_SVG)

# TODO: Better generation that doesn't require global plantuml install
%.svg: %.puml
	plantuml -svg $<

# Format everything
.PHONY: fmt
fmt:
	npx prettier . --write

# Prettier installs node_modules
node_modules: package.json package-lock.json
	npm install

# Run prettier on pre-commit
.git/hooks/pre-commit:
	cp .evertras/pre-commit.sh .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
