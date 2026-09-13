# Implementation Plan: Instance List Columns and Auth Styling

## Phase 1: Table Column Adjustments & Styling [checkpoint: dbc26e1]
- [x] Task: Update Instances Table Headers, Expansions, and Cell Rendering 504629a
    - [x] Update unit tests in `internal/run/tui/app/instance/instance_test.go` to assert the 6-column layout, absence of CONTAINERS column, and new LAST UPDATED column index
    - [x] Update `listHeaders`, `listExpansions`, and `render()` in `internal/run/tui/app/instance/instance.go` to remove CONTAINERS and reindex columns
    - [x] Ensure Auth column retains `[red]Yes` / `[green]No` colored styling matching Services
- [x] Task: Conductor - User Manual Verification 'Table Column Adjustments & Styling' (Protocol in workflow.md) dbc26e1

## Phase 2: Validation, Verification & Integration
- [x] Task: Full Suite Verification & Build 1dc6691
    - [x] Run `make test` to ensure 100% test pass rate and high test coverage
    - [x] Run `make lint` and `make build` to confirm code style and clean compilation
- [ ] Task: Conductor - User Manual Verification 'Validation, Verification & Integration' (Protocol in workflow.md)
