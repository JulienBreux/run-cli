# Implementation Plan: Place Instances before Services in Header

## Phase 1: Reorder Shortcut Registry [checkpoint: b6c03a4]
- [x] Task: Write failing unit tests for global shortcut ordering (2ca8bfc)
    - [x] Add unit test in `internal/run/tui/app/shortcut/shortcut_test.go` asserting that `ctrl+n` (Instances) appears before `ctrl+s` (Services) in `shortcut.Registry`
    - [x] Run test and confirm red failure
- [x] Task: Reorder global shortcuts in registry (40aebf1)
    - [x] Move `ctrl+n` (Instances) before `ctrl+s` (Services) in `shortcut.Registry` in `internal/run/tui/app/shortcut/shortcut.go`
    - [x] Run unit tests and confirm green success
- [x] Task: Conductor - User Manual Verification 'Phase 1: Reorder Shortcut Registry' (Protocol in workflow.md) (b6c03a4)

## Phase 2: Header Verification and Quality Assurance
- [x] Task: Verify header shortcut rendering (52d2b3d)
    - [x] Add/update test in `internal/run/tui/component/header/header_test.go` to verify header shortcut layout behavior
    - [x] Confirm tests pass
- [x] Task: Run project quality gates and documentation updates (c7a04a8)
    - [x] Run `make test` to ensure full test suite passes with >80% coverage
    - [x] Run `make lint` to verify code style
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Header Verification and Quality Assurance' (Protocol in workflow.md)
