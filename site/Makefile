.PHONY: build
build: node_modules
	REACT_APP_EVERTRAS_API_BASE_PATH="https://leaderboard-api.evertras.com" npm run build

node_modules: package.json package-lock.json
	npm install
	@touch node_modules
