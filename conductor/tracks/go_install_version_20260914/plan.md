# Implementation Plan: Resolve Version Information for go install Builds and Check for Updates

## Phase 1: Build Info Detection
- [x] Task: Write failing unit tests for runtime/debug build info detection in `pkg/version` (813d739)
    - [x] Add unit tests in `pkg/version/version_test.go` covering fallback to `runtime/debug.ReadBuildInfo()` when `Version == "dev"`
    - [x] Cover extraction of module version, commit hash (`vcs.revision`), dirty status (`vcs.modified`), and build date (`vcs.time`)
    - [x] Run tests and verify failure (Red phase)
- [ ] Task: Implement build info resolution in `pkg/version`
    - [ ] Implement detection logic using `runtime/debug.ReadBuildInfo()` with fallback precedence
    - [ ] Run unit tests and confirm success (Green phase)
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Build Info Detection' (Protocol in workflow.md)

## Phase 2: Remote Update Checker & TUI Notification
- [ ] Task: Write failing unit tests for remote GitHub release check in `pkg/version`
    - [ ] Add unit tests for latest release checking, semver comparison, and network failure tolerance
    - [ ] Run tests and verify failure (Red phase)
- [ ] Task: Implement remote release checker in `pkg/version`
    - [ ] Implement non-blocking HTTP check against GitHub Releases API with 2s timeout
    - [ ] Run unit tests and confirm success (Green phase)
- [ ] Task: Write failing unit tests for header update notification
    - [ ] Add unit test in `internal/run/tui/component/header/header_test.go` for rendering version with update notice
    - [ ] Run tests and verify failure (Red phase)
- [ ] Task: Implement asynchronous header notification in `internal/run/tui/component/header`
    - [ ] Trigger background update check on header initialization and refresh header text if an update is found
    - [ ] Run unit tests and confirm success (Green phase)
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Remote Update Checker & TUI Notification' (Protocol in workflow.md)

## Phase 3: Quality Gates & Documentation
- [ ] Task: Verify CLI version command output and integration tests
    - [ ] Update `internal/run/command/version/version_test.go` to verify version output formatting
    - [ ] Confirm all command tests pass
- [ ] Task: Update README.md and documentation
    - [ ] Document version detection and update notification feature in `README.md`
- [ ] Task: Run project quality gates
    - [ ] Run `make test` to ensure full test suite passes with >80% coverage
    - [ ] Run `make lint` to verify code style
    - [ ] Run `make build` to ensure clean build
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Quality Gates & Documentation' (Protocol in workflow.md)
