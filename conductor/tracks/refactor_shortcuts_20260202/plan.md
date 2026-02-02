# Implementation Plan - Refactor Shortcut Display

## Phase 1: Shortcut Registry Enhancements [checkpoint: 56c9be4]
- [x] Task: Update the `shortcut` package with formatting logic.
    - [x] Write unit tests for `GetByCategory` and `Format` functions.
    - [x] Implement helper to filter shortcuts by category.
    - [x] Implement `Format()` method in `shortcut` package to return the `[dodgerblue]<key> [white]Desc` string.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Shortcut Registry Enhancements' (Protocol in workflow.md) 56c9be4

## Phase 2: Footer Refactoring [checkpoint: 309f2c5]
- [x] Task: Update Service views to use the new registry. 64990 (tests pass)
    - [x] Refactor `service.Shortcuts()` and `service.DashboardShortcuts()` to use `shortcut.FormatByCategory`.
    - [x] Pass necessary overrides for dynamic states (Proxy, Auth).
- [x] Task: Update remaining functional views. 64990 (tests pass)
    - [x] Refactor `job.Shortcuts()`, `workerpool.Shortcuts()`, and `domainmapping.Shortcuts()`.
- [x] Task: Conductor - User Manual Verification 'Phase 2: Footer Refactoring' (Protocol in workflow.md) 309f2c5

## Phase 3: Header Refactoring [checkpoint: 309f2c5]
- [x] Task: Update the global header to use the registry. 65162 (tests pass)
    - [x] Refactor `internal/run/tui/component/header/header.go` to iterate through `CategoryGlobal` shortcuts.
    - [x] Remove hardcoded strings from the header's column functions.
- [x] Task: Conductor - User Manual Verification 'Phase 3: Header Refactoring' (Protocol in workflow.md) 309f2c5
