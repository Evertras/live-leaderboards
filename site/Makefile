.PHONY: build
build: generated
	rm -rf build
	npm run build

.PHONY: test
test: generated
	npm test

.PHONY: dev-run
dev-run: generated
	npm start

.PHONY: generated
generated: node_modules ./src/lib/client/api.d.ts

.PHONY: clean
clean:
	rm -rf node_modules

./src/lib/client/api.d.ts: node_modules ../specs/openapi.yaml
	npx openapi-typescript ../specs/openapi.yaml -o ./src/lib/client/api.d.ts

node_modules: package.json package-lock.json
	npm install
	@touch node_modules
