# Implementation Plan: Dynamic Terminal CLI Title

## Phase 1: Title Utility & Non-TUI Commands [checkpoint: cf526f2]
- [x] Task: Centralized Title Formatting & Terminal Control Utility b20272b
    - [x] Write unit tests for title formatting (`FormatTitle`) and terminal OSC escape sequence handling in `pkg/term`
    - [x] Implement `pkg/term` with `FormatTitle` and terminal title set/reset functions
- [x] Task: Integrate Title Management in CLI Subcommands c0d096b
    - [x] Write unit tests verifying terminal title lifecycle in Cobra subcommands
    - [x] Update `version` command and root command hooks to set terminal title and restore on exit
- [x] Task: Conductor - User Manual Verification 'Title Utility & Non-TUI Commands' (Protocol in workflow.md) cf526f2

## Phase 2: TUI Application Title Management
- [ ] Task: Integrate Dynamic Titles in TUI Navigation & Dashboard
    - [ ] Write unit tests in `internal/run/tui/app/app_test.go` verifying title updates on startup, list views, and detail dashboards
    - [ ] Implement title state management in `internal/run/tui/app/app.go` (`switchToPage`, startup, detail views)
- [ ] Task: Integrate Dynamic Titles in TUI Modals & Overlays
    - [ ] Write unit tests verifying title push and restore for modals (Help, Scale, Auth, Project, Region)
    - [ ] Implement modal title push and restore logic in `internal/run/tui/app/app.go`
- [ ] Task: Conductor - User Manual Verification 'TUI Application Title Management' (Protocol in workflow.md)

## Phase 3: Validation, Verification & Integration
- [ ] Task: Full Suite Verification & Build
    - [ ] Run `make test` across all packages to verify test coverage and pass rate
    - [ ] Run `make lint` and `make build` to confirm code style and clean compilation
- [ ] Task: Conductor - User Manual Verification 'Validation, Verification & Integration' (Protocol in workflow.md)
