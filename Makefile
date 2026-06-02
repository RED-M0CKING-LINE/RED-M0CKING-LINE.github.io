APP        := website-main-go
PKG        := ethanashley.net/$(APP)
BIN        := bin/server
WASM_OUT   := web/static/wasm/tool.wasm
WASM_EXEC  := web/static/wasm/wasm_exec.js
GO         ?= go
IMAGE      ?= $(APP):dev
CONTAINER  ?= podman # set to `docker` for Docker tooling
COMPOSE    ?= podman-compose # set to `docker compose` for Docker

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show help
	@awk 'BEGIN{FS=":.*##"; printf "Usage:\n  make <target>\n\nTargets:\n"} \
	      /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Dev
.PHONY: run
run: wasm ## Run the server locally
	$(GO) run ./cmd/server

.PHONY: fmt
fmt: ## gofmt -w on the whole module
	gofmt -w .

.PHONY: vet
vet: ## go vet ./...
	$(GO) vet ./...

.PHONY: test
test: ## go test ./...
	$(GO) test -race ./...

.PHONY: lint
lint: fmt vet ## Format + vet

# Build
.PHONY: build
build: wasm ## Build the server binary into bin/
	@mkdir -p bin
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags="-s -w" -o $(BIN) ./cmd/server

.PHONY: wasm
wasm: $(WASM_OUT) $(WASM_EXEC) ## Build the WASM tool and copy wasm_exec.js

$(WASM_OUT): wasm/tool/main.go
	@mkdir -p $(dir $@)
	GOOS=js GOARCH=wasm $(GO) build -trimpath -ldflags="-s -w" -o $@ ./wasm/tool

$(WASM_EXEC):
	@mkdir -p $(dir $@)
	cp "$$($(GO) env GOROOT)/lib/wasm/wasm_exec.js" $@

# Containers
.PHONY: podman-build
podman-build: ## Build the container image with Podman
	podman build -t $(IMAGE) -f Containerfile .

.PHONY: docker-build
docker-build: ## Build the container image with Docker
	docker build -t $(IMAGE) -f Dockerfile .

.PHONY: podman-up
podman-up: ## Bring the stack up with podman-compose
	cd deploy/compose && podman-compose up --build -d

.PHONY: podman-down
podman-down: ## Tear down podman-compose stack.
	cd deploy/compose && podman-compose down

.PHONY: docker-up
docker-up: ## Bring the stack up with docker compose
	cd deploy/compose && docker compose up --build -d

.PHONY: docker-down
docker-down: ## Tear down docker compose stack
	cd deploy/compose && docker compose down

# Utility
.PHONY: clean
clean: ## Remove build artefacts.
	rm -rf bin $(WASM_OUT) $(WASM_EXEC)

.PHONY: cookie-secret
cookie-secret: ## Print a fresh 32-byte hex secret for COOKIE_SECRET
	@openssl rand -hex 32

.PHONY: tidy
tidy: ## go mod tidy
	$(GO) mod tidy
