# Implementation Plan - Refactor Shortcut Display

## Phase 1: Shortcut Registry Enhancements [checkpoint: 56c9be4]
- [x] Task: Update the `shortcut` package with formatting logic.
    - [x] Write unit tests for `GetByCategory` and `Format` functions.
    - [x] Implement helper to filter shortcuts by category.
    - [x] Implement `Format()` method in `shortcut` package to return the `[dodgerblue]<key> [white]Desc` string.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Shortcut Registry Enhancements' (Protocol in workflow.md) 56c9be4

## Phase 2: Footer Refactoring
- [~] Task: Update Service views to use the new registry.
    - [ ] Refactor `service.Shortcuts()` and `service.DashboardShortcuts()` to use `shortcut.FormatByCategory`.
    - [ ] Pass necessary overrides for dynamic states (Proxy, Auth).
- [ ] Task: Update remaining functional views.
    - [ ] Refactor `job.Shortcuts()`, `workerpool.Shortcuts()`, and `domainmapping.Shortcuts()`.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Footer Refactoring' (Protocol in workflow.md)
