# Implementation Plan - Update Traffic Split

## Phase 1: API Client Enhancements
- [x] Task: Update the Service API client to support updating traffic splits. 061e4ab
    - [x] Write unit tests for the traffic update mapping logic and client execution.
    - [x] Implement `UpdateTraffic` method in `internal/run/api/service/service.go`.
- [x] Task: Conductor - User Manual Verification 'Phase 1: API Client Enhancements' (Protocol in workflow.md) caaadf2

## Phase 2: Traffic Split Modal (Refactor)
- [x] Task: Refactor Traffic Split Modal to use dynamic Dropdown selectors. 428e71c
    - [x] Update `Modal` signature to accept all revisions (not just selected).
    - [x] Implement dynamic form with "Add Revision" button.
    - [x] Use `pkg/dropdown` for revision selection in each row.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Traffic Split Modal (Refactor)' (Protocol in workflow.md)

## Phase 3: Integration
- [ ] Task: Connect the Service Dashboard to the new Traffic Split modal.
    - [ ] Add `t` shortcut to trigger the modal (passing full service and revision list).
    - [ ] Remove any previous selection logic if present.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Integration' (Protocol in workflow.md)