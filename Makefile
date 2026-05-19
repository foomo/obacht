.DEFAULT_GOAL:=help
-include .makerc

# --- Config -----------------------------------------------------------------

# Newline hack for error output
define br


endef

# --- Targets -----------------------------------------------------------------

# This allows us to accept extra arguments
%: .mise .lefthook
	@:

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
## Run lint & tests
check: tidy generate lint lint.rules test.race audit

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

.PHONY: lint.rules
## Lint embedded rule artifacts (shell + rego)
lint.rules: lint.shell lint.policy

.PHONY: lint.shell
## Run shellcheck on rule input scripts
lint.shell:
	@echo "〉shellcheck rules/inputs/**/*.sh"
	@shellcheck -x rules/inputs/_lib/*.sh rules/inputs/*/*.sh

.PHONY: lint.policy
## Run regal on rule policies
lint.policy:
	@echo "〉regal lint rules/policy/"
	@regal lint rules/policy/

.PHONY: fmt.rules
## Format embedded rule artifacts (shell + rego)
fmt.rules: fmt.shell fmt.policy

.PHONY: fmt.shell
## Format rule shell scripts with shfmt
fmt.shell:
	@echo "〉shfmt -w rules/inputs/"
	@shfmt -w rules/inputs/

.PHONY: fmt.policy
## Format rule rego policies with opa fmt
fmt.policy:
	@echo "〉opa fmt -w rules/policy/"
	@opa fmt -w rules/policy/

.PHONY: generate
## Run go generate
generate:
	@echo "〉go generate"
	@go generate ./...

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

.PHONY: build
## Build binary
build:
	@echo "〉go build bin/obacht"
	@rm -f bin/obacht
	@go build -o bin/obacht cmd/obacht/obacht.go

.PHONY: build.debug
## Build binary in debug mode
build.debug:
	@echo "〉go build bin/obacht (debug)"
	@rm -f bin/obacht
	@go build -gcflags "all=-N -l" -o bin/obacht cmd/obacht/obacht.go

.PHONY: install
## Run go install
install:
	@echo "〉installing obacht"
	@go install cmd/obacht/obacht.go

.PHONY: install.debug
## Run go install with debug
install.debug:
	@echo "〉installing obacht (debug)"
	@go install -gcflags "all=-N -l" cmd/obacht/obacht.go

### Security

.PHONY: audit
## Run security audit
audit:
	@echo "〉security audit"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...

### Dependencies

.PHONY: tidy
## Run go mod tidy
tidy:
	@echo "〉go mod tidy"
	@go mod tidy

.PHONY: outdated
## Show outdated direct dependencies
outdated:
	@echo "〉go mod outdated"
	@go list -u -m -json all | go-mod-outdated -update -direct

.PHONY: upgrade
## Show outdated direct dependencies
upgrade:
	@echo "〉go mod upgrade"
	@go list -u -m -f '{{if and (not .Indirect) .Update}}{{.Path}}{{end}}' all | xargs -n1 -I{} go get {}@latest
	@$(MAKE) tidy

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
# https://patorjk.com/software/taag/#p=display&f=Tmplr&t=Obacht&x=none&v=4&h=4&w=80&we=false
## Show help text
help: g=\033[0;32m
help: b=\033[0;34m
help: w=\033[0;90m
help: e=\033[0m
help:
	@echo "$(g)"
	@echo "  ┓    ┓  "
	@echo "┏┓┣┓┏┓┏┣┓╋"
	@echo "┗┛┗┛┗┻┗┛┗┗"
	@echo "with ❤ foomo by bestbytes"
	@echo "$(e)"
	@echo "$(b)Usage:$(e)\n  make [task]"
	@awk '{ \
		if($$0 ~ /^### /){ \
			if(help) printf "  %-21s $(w)%s$(e)\n\n", cmd, help; help=""; \
			printf "$(b)\n%s:$(e)\n", substr($$0,5); \
		} else if($$0 ~ /^[a-zA-Z0-9._-]+:/){ \
			cmd = substr($$0, 1, index($$0, ":")-1); \
			if(help) printf "  %-21s $(w)%s$(e)\n", cmd, help; help=""; \
		} else if($$0 ~ /^##/){ \
			help = help ? help "\n                        " substr($$0,3) : substr($$0,3); \
		} else if(help){ \
			print "\n                        $(w)" help "$(e)\n"; help=""; \
		} \
	}' $(MAKEFILE_LIST)
	@echo ""

