.PHONY: build
build: node_modules ./src/lib/client/api.d.ts
	rm -rf build
	npm run build

.PHONY: generated
generated: ./src/lib/client/api.d.ts

./src/lib/client/api.d.ts: node_modules ../specs/openapi.yaml
	npx openapi-typescript ../specs/openapi.yaml -o ./src/lib/client/api.d.ts

node_modules: package.json package-lock.json
	npm install
	@touch node_modules
