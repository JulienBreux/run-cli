# Tech Stack

## Core Technologies
- **Programming Language:** [Go (Golang)](https://go.dev/) (Version 1.26+) - The primary language for its performance, concurrency support, and strong ecosystem for CLI tools.

## Frameworks & Libraries
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra) - Used for building the CLI command structure, flag parsing, and help generation.
- **TUI Engine:** [tview](https://github.com/rivo/tview) (built on [tcell](https://github.com/gdamore/tcell)) - The foundation for the interactive terminal user interface, providing rich components like tables, forms, and layout management.

## Infrastructure & Cloud
- **Platform:** Google Cloud Platform (GCP).
- **SDKs:**
  - `cloud.google.com/go/run`: For managing Cloud Run Services, Jobs, and Revisions.
  - `cloud.google.com/go/logging`: For streaming and viewing logs.
  - `cloud.google.com/go/resourcemanager`: For project and region discovery.

## Development & Testing
- **Testing Framework:** [testify](https://github.com/stretchr/testify) - Used for assertions and mocking in unit tests.
- **Testing:** `make tests` - Runs unit tests and benchmarks.
- **Test coverage:** `make coverage-total` - Get the total number of lines covered by tests. All modules must maintain >80% code coverage.
- **Linting:** `golangci-lint` - Enforces code quality and style standards.
- **Build System:** `Makefile` - Orchestrates build, test, and linting tasks.

## Architecture
- **Standard Go Layout:**
  - `cmd/run/`: Main entry point.
  - `internal/run/`: Private implementation details including API clients, TUI logic, and data models.
  - `pkg/`: Reusable public utility packages.
