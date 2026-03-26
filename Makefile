.DEFAULT_GOAL:=help
-include .makerc

# --- Config -----------------------------------------------------------------

# Newline hack for error output
define br


endef

# --- Targets -----------------------------------------------------------------

# This allows us to accept extra arguments
%: .mise .lefthook go.work
	@:

# Ensure go.work file
go.work:
	@go work init
	@go work use -r .
	@go work sync

.PHONY: .mise
# Install dependencies
.mise:
ifeq (, $(shell command -v mise))
	$(error $(br)$(br)Please ensure you have 'mise' installed and activated!$(br)$(br)  $$ brew update$(br)  $$ brew install mise$(br)$(br)See the documentation: https://mise.jdx.dev/getting-started.html)
endif
	@mise install

.PHONY: .lefthook
# Configure git hooks for lefthook
.lefthook:
	@lefthook install --reset-hooks-path

### Tasks

.PHONY: check
## Run lint & test
check: tidy examples generate lint test.race

.PHONY: tidy
## Run go mod tidy
tidy:
	@echo "〉go mod tidy"
	go mod tidy

regal:
	@regal lint ./policies

.PHONY: lint
## Run linter
lint:
	@echo "〉golangci-lint run"
	golangci-lint run --max-same-issues 0 --max-issues-per-linter 0

.PHONY: lint.fix
## Fix lint violations
lint.fix:
	@echo "〉golangci-lint run fix"
	golangci-lint run --fix --max-same-issues 0 --max-issues-per-linter 0

.PHONY: test
## Run tests
test:
	@echo "〉go test"
	@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out -tags=safe ./...

.PHONY: test.race
## Run tests with -race
test.race:
	@echo "〉go test -race"
	@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out -tags=safe -race ./...

.PHONY: test.nocache
## Run tests with -count=1
test.nocache:
	@echo "〉go test -count=1"
	@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out -tags=safe -count=1 ./...

.PHONY: outdated
## Show outdated direct dependencies
outdated:
	@echo "〉go mod outdated"
	@go list -u -m -json all | go-mod-outdated -update -direct

.PHONY: build
## Build binary
build:
	@echo "〉go build bin/bouncer"
	@rm -f bin/bouncer
	@go build -o bin/bouncer cmd/bouncer/main.go

.PHONY: build.debug
## Build binary in debug mode
build.debug:
	@echo "〉go build bin/bouncer (debug)"
	@rm -f bin/bouncer
	@go build -gcflags "all=-N -l" -o bin/bouncer cmd/bouncer/main.go

.PHONY: install
## Run go install
install:
	@echo "〉installing bouncer"
	@go install cmd/bouncer/main.go

.PHONY: install.debug
## Run go install with debug
install.debug:
	@echo "〉installing bouncer (debug)"
	@go install -gcflags "all=-N -l" cmd/bouncer/main.go

.PHONY: generate
## Run go generate
generate:
	@echo "〉go generate"
	@go generate ./...

### Documentation

.PHONY: docs
## Open docs
docs:
	@echo "〉starting docs"
	@cd docs && bun install && bun run dev

.PHONY: docs.build
## Open docs
docs.build:
	@echo "〉building docs"
	@cd docs && bun install && bun run build

.PHONY: godocs
## Open go docs
godocs:
	@echo "〉starting go docs"
	@go doc -http

### Utils

.PHONY: help
## Show help text
help:
	@echo "bouncer\n"
	@echo "Usage:\n  make [task]"
	@awk '{ \
		if($$0 ~ /^### /){ \
			if(help) printf "%-23s %s\n\n", cmd, help; help=""; \
			printf "\n%s:\n", substr($$0,5); \
		} else if($$0 ~ /^[a-zA-Z0-9._-]+:/){ \
			cmd = substr($$0, 1, index($$0, ":")-1); \
			if(help) printf "  %-23s %s\n", cmd, help; help=""; \
		} else if($$0 ~ /^##/){ \
			help = help ? help "\n                        " substr($$0,3) : substr($$0,3); \
		} else if(help){ \
			print "\n                        " help "\n"; help=""; \
		} \
	}' $(MAKEFILE_LIST)
