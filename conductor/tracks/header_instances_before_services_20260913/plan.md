# Implementation Plan: Place Instances before Services in Header

## Phase 1: Reorder Shortcut Registry
- [x] Task: Write failing unit tests for global shortcut ordering (2ca8bfc)
    - [x] Add unit test in `internal/run/tui/app/shortcut/shortcut_test.go` asserting that `ctrl+n` (Instances) appears before `ctrl+s` (Services) in `shortcut.Registry`
    - [x] Run test and confirm red failure
- [x] Task: Reorder global shortcuts in registry (40aebf1)
    - [x] Move `ctrl+n` (Instances) before `ctrl+s` (Services) in `shortcut.Registry` in `internal/run/tui/app/shortcut/shortcut.go`
    - [x] Run unit tests and confirm green success
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Reorder Shortcut Registry' (Protocol in workflow.md)

## Phase 2: Header Verification and Quality Assurance
- [ ] Task: Verify header shortcut rendering
    - [ ] Add/update test in `internal/run/tui/component/header/header_test.go` to verify header shortcut layout behavior
    - [ ] Confirm tests pass
- [ ] Task: Run project quality gates and documentation updates
    - [ ] Run `make test` to ensure full test suite passes with >80% coverage
    - [ ] Run `make lint` to verify code style
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Header Verification and Quality Assurance' (Protocol in workflow.md)
