[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/obacht?style=flat-square)](https://goreportcard.com/report/github.com/foomo/obacht)
[![GoDoc](https://img.shields.io/badge/GoDoc-✓-informational.svg?style=flat-square&logo=go)](https://godoc.org/github.com/foomo/obacht)
[![GitHub Downloads](https://img.shields.io/github/downloads/foomo/obacht/total.svg?style=flat-square&logo=github)](https://github.com/foomo/obacht/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/foomo/obacht.svg?style=flat-square&logo=docker)](https://hub.docker.com/r/foomo/obacht)
[![GitHub Stars](https://img.shields.io/github/stars/foomo/obacht.svg?style=flat-square&logo=github)](https://github.com/foomo/obacht)

<p align="center">
  <img alt="obacht" src="docs/public/logo.png" width="400" height="400"/>
</p>

# obacht

> Security scanner for developer environments

obacht inspects your local development setup for security misconfigurations — insecure file permissions, exposed credentials, weak SSH/Git settings, risky Docker access — using an embedded [OPA](https://www.openpolicyagent.org/) engine and Rego policies. It is lightweight, read-only, and requires no agent or endpoint management platform.

## Features

- **98 built-in rules** across 12 categories: SSH, Git, Docker, Kubernetes, env, shell, tools, PATH, OS, credentials, privacy
- **OPA-powered** with an embedded Rego engine — no external dependencies
- **Read-only collectors** — never modifies system state
- **Extensible** via `--rules-dir` for custom Rego policies
- **Pretty TUI** or machine-readable JSON output for CI

## Installation

<details>
<summary><b>Homebrew</b> (macOS / Linux)</summary>

```shell
brew install foomo/tap/obacht
```

See the [foomo/homebrew-tap](https://github.com/foomo/homebrew-tap) repository.

</details>

<details>
<summary><b>Docker</b></summary>

```shell
docker run --rm foomo/obacht:latest scan
```

Multi-arch images (`amd64`, `arm64`) are published to [Docker Hub](https://hub.docker.com/r/foomo/obacht).

</details>

<details>
<summary><b>mise</b></summary>

```shell
mise use github:foomo/obacht
```

or run directly:

```shell
mise x github:foomo/obacht -- scan
```

See [mise.jdx.dev](https://mise.jdx.dev).

</details>

<details>
<summary><b>Binary release</b></summary>

Download the archive for your OS/arch from the [releases page](https://github.com/foomo/obacht/releases) and extract `obacht` into your `$PATH`.

</details>

<details>
<summary><b>go install</b></summary>

```shell
go install github.com/foomo/obacht/cmd/obacht@latest
```

Requires Go 1.26+.

</details>

## Usage

```shell
$ obacht --help
Security configuration scanner for developer environments

Usage:
  obacht [flags]
  obacht [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  doctor      Check obacht dependencies and configuration
  explain     Show detailed information about a rule
  help        Help about any command
  scan        Scan the local development environment for security issues
  version     Print version information

Flags:
      --format string      output format (pretty, json) (default "pretty")
  -h, --help               help for obacht
      --rules-dir string   use rules from this directory instead of embedded rules
      --verbose            enable verbose output
      --version            print version information

Use "obacht [command] --help" for more information about a command.
```

## Resources

- [Foomo Security](https://www.foomo.org/blog/tag/security/)
- [Pareto Security](https://github.com/ParetoSecurity/pareto-mac)

## How to Contribute

Contributions are welcome! Please read the [contributing guide](CONTRIBUTING.md).

![Contributors](https://contributors-table.vercel.app/image?repo=foomo/obacht&width=50&columns=15)

## License

Distributed under MIT License, please see license file within the code for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_
