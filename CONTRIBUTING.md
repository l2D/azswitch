# Contributing to azswitch

Thank you for your interest in contributing to azswitch! This document provides guidelines and information about contributing.

## Development Setup

### Prerequisites

- Go 1.25 or later
- Azure CLI (for testing)
- golangci-lint (for linting)
- Docker (optional, for container builds)

### Getting Started

1. Fork the repository
2. Clone your fork:

   ```bash
   git clone https://github.com/YOUR_USERNAME/azswitch.git
   cd azswitch
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Build:

   ```bash
   make build
   ```

5. Run tests:

   ```bash
   make test
   ```

## Making Changes

### Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Run `make lint` to check for issues
- Write tests for new functionality

### Commit Messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Test changes
- `chore:` - Maintenance tasks

Example:

```sh
feat(tui): add search filtering for subscriptions

- Add fuzzy search to subscription list
- Highlight matching characters
```

### Pull Requests

1. Create a feature branch from `main`
2. Make your changes
3. Run tests and linting
4. Push your branch
5. Open a Pull Request

## Project Structure

```sh
azswitch/
├── cmd/azswitch/       # Application entry point
├── internal/
│   ├── azure/          # Azure CLI wrapper
│   ├── tui/            # Bubble Tea TUI
│   └── version/        # Version info
├── .github/workflows/  # CI/CD
└── Makefile           # Build automation
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage
```

### Mock Client

Use `azure.NewMockClient()` for testing TUI components without Azure CLI.

## Questions?

Feel free to open an issue for questions or discussions.
