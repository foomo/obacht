# Getting Started

## Installation

### Using mise (recommended)

```bash
mise install obacht
```

### Using Homebrew

```bash
brew install franklinkim/tap/obacht
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
