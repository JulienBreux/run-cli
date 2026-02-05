# Specification: Increase Test Coverage to >80%

## Overview
The current code coverage of the `run-cli` project is approximately 72%. To ensure high reliability and meet project standards, this track aims to increase the overall test coverage to at least 80% by adding comprehensive unit tests across key components.

## Functional Requirements
- Identify all packages and files with coverage below 80% using existing tooling.
- Implement unit tests for core models in `internal/run/model` to ensure business logic is fully verified.
- Enhance test suites for API clients in `internal/run/api`, ensuring edge cases and error conditions are covered.
- Expand coverage for CLI commands in `internal/run/command` and TUI components in `internal/run/tui`.
- Utilize `testify/mock` for mocking Google Cloud APIs and other external dependencies.
- Implement file-based fixtures for complex API response simulations where appropriate.

## Non-Functional Requirements
- **Performance:** Test execution time should remain efficient to maintain a fast development cycle.
- **Maintainability:** New tests must follow existing project conventions and be easy to understand and update.

## Acceptance Criteria
- `make test` reports an overall project code coverage of >80%.
- `codecov.yml` report confirms the coverage increase.
- All new and existing tests pass successfully.
- No regression in existing functionality.

## Out of Scope
- Integration tests requiring actual GCP credentials or resources.
