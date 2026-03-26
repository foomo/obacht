# Getting Started

## Installation

### Using mise (recommended)

```bash
mise install bouncer
```

### Using Homebrew

```bash
brew install franklinkim/tap/bouncer
```

### From Source

```bash
go install github.com/franklinkim/bouncer/cmd/bouncer@latest
```

## Prerequisites

Bouncer requires the [Open Policy Agent (OPA)](https://www.openpolicyagent.org/) binary:

```bash
brew install opa
# or
mise install opa
```

Verify your setup:

```bash
bouncer doctor
```

## First Scan

Run a full scan of your development environment:

```bash
bouncer scan
```

For JSON output:

```bash
bouncer scan --format json
```

Filter by category:

```bash
bouncer scan --category ssh,git
```
