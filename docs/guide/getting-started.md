# Getting Started

## Installation

### Using Homebrew

```bash
brew install foomo/tap/obacht
```

### Using mise

```bash
mise install github:foomo/obacht
```

### From Source

```bash
go install github.com/foomo/obacht/cmd/obacht@latest
```

Verify your setup:

```bash
obacht doctor
```

## First Scan

Run a full scan of your development environment:

```bash
obacht scan
```

For JSON output:

```bash
obacht scan --format json
```

Filter by category:

```bash
obacht scan --category ssh,git
```
