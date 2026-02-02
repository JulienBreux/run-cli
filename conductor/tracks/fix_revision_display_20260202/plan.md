# Implementation Plan - Fix Stale Revision Display

## Phase 1: Component Reset Logic [checkpoint: d469a67]
- [x] Task: Implement Clear methods in Revision TUI components. d469a67
    - [x] Write unit tests for `ListComponent.Clear()` and `DetailComponent.Clear()` in `internal/run/tui/app/service/revision/revision_test.go`.
    - [x] Implement `Clear()` in `internal/run/tui/app/service/revision/revision.go` to wipe internal state and reset UI text/cells.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Component Reset Logic' (Protocol in workflow.md) d469a67

## Phase 2: Dashboard State Synchronization [checkpoint: 1b2696f]
- [x] Task: Integrate data clearing into the Dashboard lifecycle. 1b2696f
    - [x] Write unit tests to verify `DashboardReload` calls clearing logic.
    - [x] Update `DashboardReload` in `internal/run/tui/app/service/dashboard.go` to clear components before triggering async data fetch.
    - [x] Update navigation handlers in `internal/run/tui/app/app.go` to clear dashboard data when exiting the Service Dashboard view.
- [x] Task: Conductor - User Manual Verification 'Phase 2: Dashboard State Synchronization' (Protocol in workflow.md) 1b2696f
