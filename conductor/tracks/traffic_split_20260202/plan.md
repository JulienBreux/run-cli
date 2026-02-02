# Implementation Plan - Update Traffic Split

## Phase 1: API Client Enhancements [checkpoint: caaadf2]
- [x] Task: Update the Service API client to support updating traffic splits. 061e4ab
    - [x] Write unit tests for the traffic update mapping logic and client execution.
    - [x] Implement `UpdateTraffic` method in `internal/run/api/service/service.go`.
- [x] Task: Conductor - User Manual Verification 'Phase 1: API Client Enhancements' (Protocol in workflow.md) caaadf2

## Phase 2: Traffic Split Modal [checkpoint: 0870cff]
- [x] Task: Create a new TUI modal for entering traffic split percentages. a855d3c
    - [x] Write unit tests for percentage validation logic (must sum to 100).
    - [x] Implement the modal component in `internal/run/tui/app/service/traffic/split.go`.
    - [x] Ensure the modal dynamically lists selected revisions and provides input fields.
- [x] Task: Conductor - User Manual Verification 'Phase 2: Traffic Split Modal' (Protocol in workflow.md) 0870cff

## Phase 3: Integration and Selection
- [x] Task: Implement multi-selection logic for the Revisions table. b29dbbe
    - [x] Write unit tests for selection tracking logic.
    - [x] Update `table` component or `dashboard.go` to support marking/selecting multiple rows.
- [x] Task: Connect the Revisions table selection to the Traffic Split modal. b29dbbe
    - [x] Add a shortcut (e.g., `t`) to trigger the modal for selected revisions.
    - [x] Integrate the API call with the modal's "Save" action using a spinner for feedback.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Integration and Selection' (Protocol in workflow.md)
