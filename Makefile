RUN_CONFIG_FILE_LOCAL=run.yaml.dist
RUN_CONFIG_FILE=/etc/run.yaml
RUN_BIN_DIR=./bin
RUN_BIN=./bin/run
RUN_ASSETS_DIR=./docs/assets
RUN_DEMO_CAST_FILE=run.cast
RUN_DEMO_GIF_FILE=run.cast

generate: ## Run go generate
	go generate

lint: ## Lint code
	golangci-lint run

test: ## Test packages
	go test -count=1 -failfast -cover -coverprofile=coverage.txt -v ./...

coverage: ## Test coverage with default output
	go tool cover -func=coverage.txt

coverage-html: ## Test coverage with html output
	go tool cover -html=coverage.html

clean: ## Clean project
	rm -Rf ${RUN_BIN_DIR}
	rm -Rf coverage.txt

build: clean ## Build local binary
	mkdir -p ${RUN_BIN_DIR}
	go build -o ${RUN_BIN_DIR} ./cmd/run

build-image: ## Build local image
	docker build -t ghcr.io/julienbreux/run:latest .

run: build ## Run local binary
	${RUN_BIN}

run-container: ## Run prepared local container
	docker run --rm -v $(PWD)/${RUN_CONFIG_FILE_LOCAL}:${RUN_CONFIG_FILE} -e RUN_CONFIG_FILE=${RUN_CONFIG_FILE} julienbreux/run:latest

demo-record: build ## Demo to .cast
	asciinema rec --command ${RUN_BIN} --title "Cloud Run CLI Demo" ${RUN_ASSETS_DIR}/${RUN_DEMO_CAST_FILE} --overwrite

demo-to-gif: ## Demo to gif
	agg ${RUN_ASSETS_DIR}/${RUN_DEMO_CAST_FILE} ${RUN_ASSETS_DIR}/${RUN_DEMO_GIF_FILE}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate lint test coverage coverage-html clean build build-image run run-container demo-record demo-to-gif help
