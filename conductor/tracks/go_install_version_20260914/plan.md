# Implementation Plan: Resolve Version Information for go install Builds and Check for Updates

## Phase 1: Build Info Detection [checkpoint: b63a49a]
- [x] Task: Write failing unit tests for runtime/debug build info detection in `pkg/version` (813d739)
    - [x] Add unit tests in `pkg/version/version_test.go` covering fallback to `runtime/debug.ReadBuildInfo()` when `Version == "dev"`
    - [x] Cover extraction of module version, commit hash (`vcs.revision`), dirty status (`vcs.modified`), and build date (`vcs.time`)
    - [x] Run tests and verify failure (Red phase)
- [x] Task: Implement build info resolution in `pkg/version` (5dd25e2)
    - [x] Implement detection logic using `runtime/debug.ReadBuildInfo()` with fallback precedence
    - [x] Run unit tests and confirm success (Green phase)
- [x] Task: Conductor - User Manual Verification 'Phase 1: Build Info Detection' (Protocol in workflow.md) (b63a49a)

## Phase 2: Remote Update Checker & TUI Notification
- [x] Task: Write failing unit tests for remote GitHub release check in `pkg/version` (198671d)
    - [x] Add unit tests for latest release checking, semver comparison, and network failure tolerance
    - [x] Run tests and verify failure (Red phase)
- [x] Task: Implement remote release checker in `pkg/version` (12a85db)
    - [x] Implement non-blocking HTTP check against GitHub Releases API with 2s timeout
    - [x] Run unit tests and confirm success (Green phase)
- [x] Task: Write failing unit tests for header update notification (0ce7321)
    - [x] Add unit test in `internal/run/tui/component/header/header_test.go` for rendering version with update notice
    - [x] Run tests and verify failure (Red phase)
- [x] Task: Implement asynchronous header notification in `internal/run/tui/component/header` (ff81d33)
    - [x] Trigger background update check on header initialization and refresh header text if an update is found
    - [x] Run unit tests and confirm success (Green phase)
- [x] Task: Conductor - User Manual Verification 'Phase 2: Remote Update Checker & TUI Notification' (Protocol in workflow.md) [checkpoint: edfce46]

## Phase 3: Quality Gates & Documentation
- [x] Task: Verify CLI version command output and integration tests (79a285c)
    - [x] Update `internal/run/command/version/version_test.go` to verify version output formatting
    - [x] Confirm all command tests pass
- [x] Task: Update README.md and documentation (2c96526)
    - [x] Document version detection and update notification feature in `README.md`
- [x] Task: Run project quality gates
    - [x] Run `make test` to ensure full test suite passes with >80% coverage
    - [x] Run `make lint` to verify code style
    - [x] Run `make build` to ensure clean build
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Quality Gates & Documentation' (Protocol in workflow.md)
