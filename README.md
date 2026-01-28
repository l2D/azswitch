# azswitch

A TUI application for switching Azure tenants, directories, and subscriptions.

![Demo](docs/demo.gif)

## Features

- **Interactive TUI** - Navigate with keyboard (vim-style j/k or arrows)
- **View Current Account** - See active user, tenant, and subscription
- **Switch Subscriptions** - Quick selection from available subscriptions
- **Switch Tenants** - Re-authenticate to a different Azure AD tenant
- **CLI Mode** - Non-interactive flags for scripting

## Installation

### Homebrew (macOS/Linux)

```bash
brew install l2D/tap/azswitch
```

### Go Install

```bash
go install github.com/l2D/azswitch/cmd/azswitch@latest
```

### Binary Download

Download the latest release from [GitHub Releases](https://github.com/l2D/azswitch/releases).

### Docker

```bash
docker run --rm -it -v ~/.azure:/root/.azure ghcr.io/l2d/azswitch
```

## Prerequisites

- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) must be installed
- Must be logged in (`az login`)

## Usage

### Interactive Mode (Default)

```bash
azswitch
```

### CLI Flags

```bash
# Show current account
azswitch --current

# List all subscriptions
azswitch --list

# Switch to subscription by name or ID
azswitch --subscription "My Subscription"
azswitch --subscription xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

# Switch to a different tenant
azswitch --tenant xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

## Key Bindings

| Key | Action |
|-----|--------|
| `j` / `Down` | Move cursor down |
| `k` / `Up` | Move cursor up |
| `Enter` | Select item |
| `Tab` | Switch between subscriptions/tenants view |
| `?` | Toggle help |
| `q` / `Ctrl+C` | Quit |

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Lint

```bash
make lint
```

## License

[MIT](LICENSE)
