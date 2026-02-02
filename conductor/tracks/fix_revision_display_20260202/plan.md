# Implementation Plan - Fix Stale Revision Display

## Phase 1: Component Reset Logic [checkpoint: d469a67]
- [x] Task: Implement Clear methods in Revision TUI components. d469a67
    - [x] Write unit tests for `ListComponent.Clear()` and `DetailComponent.Clear()` in `internal/run/tui/app/service/revision/revision_test.go`.
    - [x] Implement `Clear()` in `internal/run/tui/app/service/revision/revision.go` to wipe internal state and reset UI text/cells.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Component Reset Logic' (Protocol in workflow.md) d469a67

## Phase 2: Dashboard State Synchronization
- [~] Task: Integrate data clearing into the Dashboard lifecycle.
