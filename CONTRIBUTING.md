# Contributing to Run CLI

Thank you for your interest in contributing to **Run CLI**! We welcome contributions of all kinds: bug reports, documentation updates, feature requests, and code contributions.

Please review this document to ensure a smooth contribution process.

---

## Code of Conduct

This project and everyone participating in it is governed by the [Run CLI Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior following the guidelines in [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

---

## How to Contribute

### Reporting Bugs & Feature Requests

- Before opening a new issue, please search existing [GitHub Issues](https://github.com/JulienBreux/run-cli/issues) to avoid duplicates.
- To report a bug, open an issue using the bug report template, including reproducible steps, environment details, and expected vs. actual behavior.
- To request a feature or enhancement, open an issue explaining the use case and desired functionality.
- For usage questions or general discussions, see [SUPPORT.md](SUPPORT.md).

---

## Development Setup

### Prerequisites

- **Go**: Version 1.27 or higher.
- **Make**: For running project automation targets.
- **Docker**: Optional, for container-based builds and tests.
- **golangci-lint**: For static code analysis and linting.
- **Google Cloud SDK (`gcloud`)**: Configured with appropriate Cloud Run permissions for end-to-end testing against GCP resources.

### Getting the Code

1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/<your-username>/run-cli.git
   cd run-cli
   ```
3. Set upstream remote:
   ```bash
   git remote add upstream https://github.com/JulienBreux/run-cli.git
   ```

---

## Development Workflow

The repository includes a `Makefile` with targets for standard development tasks:

- **Build**: Compile the binary to `./bin/run`:
  ```bash
  make build
  ```
  *(Note: Always use `make build` rather than direct `go build` to ensure proper build flags and paths.)*

- **Run**: Build and run the local CLI binary:
  ```bash
  make run
  ```

- **Test**: Run unit test suites with coverage analysis:
  ```bash
  make test
  ```

- **Lint**: Run static analysis (`golangci-lint`):
  ```bash
  make lint
  ```

- **Docker**:
  - Build container image: `make build-image`
  - Run container: `make run-container`

---

## Coding Guidelines and Conventions

- **Project Structure**:
  - `cmd/run/`: Application entry point (`main.go`).
  - `internal/run/command/`: CLI command definitions and Cobra flags.
  - `internal/run/tui/`: Terminal User Interface implementation using `tview` and `tcell`.
    - `app/`: Main dashboard, resource pages (`service`, `job`, `worker`, `instance`), and navigation.
    - `component/`: Reusable TUI components (header, table, spinner, loader, modals).
  - `internal/run/model/`: Domain data models and API wrappers.
  - `pkg/`: Public reusable packages.
- **Unit Tests Required**: Always add unit tests for every code change or new feature. We follow Test-Driven Development (TDD) principles and target >80% test coverage for new code.
- **License Headers**: Every new Go source file must contain the project's standard Apache 2.0 license header.
- **Commit Messages**: Follow the [Conventional Commits](https://www.conventionalcommits.org/) format:
  ```text
  <type>(<scope>): <short description>
  ```
  Common types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`.

---

## Submitting Pull Requests

1. **Create a branch**:
   ```bash
   git checkout -b feature/my-new-feature
   # or
   git checkout -b fix/issue-description
   ```
2. **Make your changes**: Keep commits focused and atomic.
3. **Verify quality gates**:
   - Run `make test` — all tests must pass.
   - Run `make lint` — ensure no linting warnings or errors.
4. **Push your branch**:
   ```bash
   git push origin feature/my-new-feature
   ```
5. **Open a Pull Request**:
   - Provide a clear PR title and fill out the [Pull Request Template](.github/pull_request_template.md).
   - Reference any related issues (e.g., `Fixes #51`).

---

## Community & Maintainers

- **Maintainers & Contributors**: See [MAINTAINERS.md](MAINTAINERS.md) for the list of project maintainers, roles, and community contributors.
- **Security**: To report a security vulnerability, please refer to [SECURITY.md](SECURITY.md).
- **Support**: For troubleshooting and assistance, refer to [SUPPORT.md](SUPPORT.md).
