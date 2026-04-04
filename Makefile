.DEFAULT_GOAL := help

VERSION := $(shell cat server.go | grep "const Version ="| cut -d"\"" -f2)
GOFLAGS=-mod=mod

# === Tool Versions (pinned) ===
GOLANGCI_VERSION := 2.11.4
GQLGEN_VERSION := v0.17.86
HADOLINT_VERSION := 2.14.0
ACT_VERSION := 0.2.87
NVM_VERSION := 0.40.4

# Parse Go version from go.mod
GO_VERSION := $(shell grep -oP '^go \K[0-9.]+' go.mod)

# Helper: run a command under the correct Go version
# In CI, actions/setup-go provides Go directly — gvm is not needed.
# Locally, gvm sets GOROOT/GOPATH/PATH in a subshell.
HAS_GVM := $(shell [ -s "$$HOME/.gvm/scripts/gvm" ] && echo true || echo false)
GVM_SHA := dd6525539fa4b771840846f8319fad303c7d0a8d2

define go-exec
$(if $(filter true,$(HAS_GVM)),bash -c '. $$GVM_ROOT/scripts/gvm && gvm use go$(GO_VERSION) >/dev/null && $(1)',bash -c '$(1)')
endef

#help: @ List available tasks
help:
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-22s\033[0m - %s\n", $$1, $$2}'

#clean: @ Cleanup
clean:
	@rm -rf ./.bin/
	@rm -rf vendor/
	@mkdir ./.bin/

#generate: @ Generate GraphQL go source code
generate: deps
	@rm -rf graph/model
	@rm -rf graph/generated
	@$(call go-exec,export GOFLAGS=$(GOFLAGS) && go run github.com/99designs/gqlgen generate)

#test: @ Run tests
test: generate
	@$(call go-exec,export GOFLAGS=$(GOFLAGS) && go test -v ./...)

#build: @ Build GraphQL API
build: generate
	@$(call go-exec,export GOFLAGS=$(GOFLAGS) && go build -o ./.bin/server server.go)

#run: @ Run GraphQL API
run: build kill-backend
	@$(call go-exec,export GOFLAGS=$(GOFLAGS) && go run server.go)

#image-build: @ Build Docker image
image-build: generate
	@docker buildx build --load -t gqlgen-graphql-subscriptions .

#build-frontend: @ Build JS client frontend
build-frontend:
	@cd ./frontend && yarn install && yarn build

#run-frontend: @ Run JS client frontend
run-frontend: build-frontend
	@cd ./frontend && yarn start

#image-frontend: @ Build JS client Docker image
image-frontend: build-frontend
	@cd ./frontend && docker build -t gqlgen-graphql-frontend .

#get: @ Download and install packages
get: clean
	@$(call go-exec,export GOFLAGS=$(GOFLAGS) && go get . && go mod tidy)

#deps: @ Install required tools (idempotent)
deps:
	@if [ -z "$$CI" ] && [ ! -s "$$HOME/.gvm/scripts/gvm" ]; then \
		echo "Installing gvm (Go Version Manager)..."; \
		curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/$(GVM_SHA)/binscripts/gvm-installer | bash -s $(GVM_SHA); \
		echo ""; \
		echo "gvm installed. Please restart your shell or run:"; \
		echo "  source $$HOME/.gvm/scripts/gvm"; \
		echo "Then re-run 'make deps' to install Go $(GO_VERSION) via gvm."; \
		exit 0; \
	fi
	@if [ "$(HAS_GVM)" = "true" ]; then \
		bash -c '. $$GVM_ROOT/scripts/gvm && gvm list' 2>/dev/null | grep -q "go$(GO_VERSION)" || { \
			echo "Installing Go $(GO_VERSION) via gvm..."; \
			bash -c '. $$GVM_ROOT/scripts/gvm && gvm install go$(GO_VERSION) -B'; \
		}; \
	else \
		command -v go >/dev/null 2>&1 || { echo "Error: Go required. Install gvm from https://github.com/moovweb/gvm or Go from https://go.dev/dl/"; exit 1; }; \
	fi
	@$(call go-exec,command -v golangci-lint) >/dev/null 2>&1 || { echo "Installing golangci-lint $(GOLANGCI_VERSION)..."; \
		$(call go-exec,go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(GOLANGCI_VERSION)); \
	}
	@$(call go-exec,command -v gqlgen) >/dev/null 2>&1 || { echo "Installing gqlgen $(GQLGEN_VERSION)..."; \
		$(call go-exec,export GOFLAGS=$(GOFLAGS) && go install github.com/99designs/gqlgen@$(GQLGEN_VERSION)); \
	}
	@command -v yarn >/dev/null 2>&1 || { echo "Installing yarn..."; npm install -g yarn; }

#deps-act: @ Install act for local CI (idempotent)
deps-act:
	@command -v act >/dev/null 2>&1 || { echo "Installing act $(ACT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash -s -- -b /usr/local/bin v$(ACT_VERSION); \
	}

#deps-hadolint: @ Install hadolint for Dockerfile linting
deps-hadolint:
	@command -v hadolint >/dev/null 2>&1 || { echo "Installing hadolint $(HADOLINT_VERSION)..."; \
		curl -sSfL -o /tmp/hadolint https://github.com/hadolint/hadolint/releases/download/v$(HADOLINT_VERSION)/hadolint-Linux-x86_64 && \
		install -m 755 /tmp/hadolint /usr/local/bin/hadolint && \
		rm -f /tmp/hadolint; \
	}

#lint: @ Run golangci-lint (includes gocritic) and hadolint
lint: deps deps-hadolint
	@$(call go-exec,golangci-lint run ./...)
	@hadolint Dockerfile
	@hadolint frontend/Dockerfile

#ci: @ Run full local CI pipeline
ci: deps generate lint test build
	@echo "Local CI pipeline passed."

#ci-run: @ Run GitHub Actions workflow locally via act
ci-run: deps-act
	@act push --container-architecture linux/amd64 -W .github/workflows/ci.yml

#release: @ Create and push a new tag
release:
	@bash -c 'read -p "New tag (current: $(VERSION)): " newtag && \
		echo "$$newtag" | grep -qE "^v[0-9]+\.[0-9]+\.[0-9]+$$" || { echo "Error: Tag must match vN.N.N"; exit 1; } && \
		echo -n "Create and push $$newtag? [y/N] " && read ans && [ "$${ans:-N}" = y ] && \
		sed -i "s/const Version = \"$(VERSION)\"/const Version = \"$$newtag\"/" server.go && \
		git add -A && \
		git commit -s -m "Cut $$newtag release" && \
		git tag $$newtag && \
		git push origin $$newtag && \
		git push && \
		echo "Done."'

#update: @ Update dependencies to latest versions
update: clean
	@$(call go-exec,export GOFLAGS=$(GOFLAGS) && go get -u && go mod tidy)

#version: @ Print current version(tag)
version:
	@echo ${VERSION}

#redis-up: @ Start Redis
redis-up: redis-down
	@docker compose up

#redis-down: @ Stop Redis
redis-down:
	@docker compose down -v --remove-orphans

#kill-backend: @ Kill all backend server processes and free port 8080
kill-backend:
	@ps aux | grep -E "(\.bin/server|server\.go)" | grep -v grep | awk '{print $$2}' | xargs -r kill -9 && echo "Killed old server processes" || echo "No server processes found"
	@sleep 1
	@lsof -i:8080 2>/dev/null || echo "Port 8080 is now free"

#renovate-bootstrap: @ Install nvm and npm for Renovate
renovate-bootstrap:
	@command -v node >/dev/null 2>&1 || { \
		echo "Installing nvm $(NVM_VERSION)..."; \
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v$(NVM_VERSION)/install.sh | bash; \
		export NVM_DIR="$$HOME/.nvm"; \
		[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh"; \
		nvm install --lts; \
	}

#renovate-validate: @ Validate Renovate configuration
renovate-validate: renovate-bootstrap
	@npx --yes renovate --platform=local

.PHONY: help clean generate test build run image-build \
	build-frontend run-frontend image-frontend \
	get deps deps-act deps-hadolint lint ci ci-run release update version \
	redis-up redis-down kill-backend \
	renovate-bootstrap renovate-validate
