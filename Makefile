RUN_CONFIG_FILE_LOCAL=run.yaml.dist
RUN_CONFIG_FILE=/etc/run.yaml
RUN_BIN_DIR=./bin
RUN_BIN=./bin/run
RUN_ASSETS_DIR=./docs/assets
RUN_DEMO_CAST_FILE=run.cast
RUN_DEMO_GIF_FILE=run.cast

.DEFAULT_GOAL := help

generate: ## Run go generate
	go generate

lint: ## Lint code
	golangci-lint run

test: ## Test packages
	go test -count=1 -failfast -cover -coverprofile=coverage.txt -v ./...

coverage: test ## Test coverage with default output
	go tool cover -func=coverage.txt

coverage-total: coverage # Get the total number of lines covered by tests
	go tool cover -func=coverage.txt | fgrep total | awk '{print substr($$3, 1, length($$3)-1)}'

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

regions-update: ## Update list of regions in region.go using gcloud
	@which gcloud > /dev/null 2>&1 || (echo "Error: gcloud CLI is required but not installed." && exit 1)
	@echo "Fetching regions from gcloud..."
	@REGIONS=$$(gcloud compute regions list | tail -n +2 | awk '{print $$1}' | sort -u); \
	if [ -z "$$REGIONS" ]; then \
		echo "Error: No regions returned from gcloud."; \
		exit 1; \
	fi; \
	FORMATTED_REGIONS=$$(echo "$$REGIONS" | sed 's/.*/\t\t"&",/'); \
	printf '/*\nCopyright 2026 Julien Breux\n\nLicensed under the Apache License, Version 2.0 (the "License");\nyou may not use this file except in compliance with the License.\nYou may obtain a copy of the License at\n\n    https://www.apache.org/licenses/LICENSE-2.0\n\nUnless required by applicable law or agreed to in writing, software\ndistributed under the License is distributed on an "AS IS" BASIS,\nWITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\nSee the License for the specific language governing permissions and\nlimitations under the License.\n*/\n\npackage region\n\n// Represents all regions.\nconst ALL = "all"\n\n// List returns a list of supported Cloud Run regions.\nfunc List() []string {\n\treturn []string{\n%s\n\t}\n}\n' "$$FORMATTED_REGIONS" > internal/run/api/region/region.go
	@go fmt ./internal/run/api/region/...
	@echo "Regions successfully updated in internal/run/api/region/region.go"

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate lint test coverage coverage-total coverage-html clean build build-image run run-container demo-record demo-to-gif help regions-update
