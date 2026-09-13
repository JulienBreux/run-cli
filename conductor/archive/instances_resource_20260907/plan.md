# Implementation Plan: Cloud Run Instances Resource Support

## Phase 1: Instances Data Models & API Client
- [x] Task: Create Instance Data Models (TDD) [a7aebfe]
    - [x] Write unit tests for instance model and converters in `internal/run/model/instance/instance_test.go`
    - [x] Implement instance data structures, status helpers, and protobuf converters in `internal/run/model/instance/instance.go`
- [x] Task: Create Instances API Client (TDD) [e07e046]
    - [x] Write unit tests and mocks for `InstancesClientWrapper` in `internal/run/api/instance/client_test.go`
    - [x] Implement `InstancesClientWrapper` and API methods (`ListInstances`, `GetInstance`, `StartInstance`, `StopInstance`, `DeleteInstance`) in `internal/run/api/instance/client.go`
- [x] Task: Conductor - User Manual Verification 'Instances Data Models & API Client' (Protocol in workflow.md)

## Phase 2: TUI Components for Instances (List, Dashboard, Lifecycle & Modals) [checkpoint: 8085198]
- [x] Task: Implement Instances List View (TDD) [d579b36]
    - [x] Write unit tests for Instances table rendering and key bindings in `internal/run/tui/app/instance/instance_test.go`
    - [x] Implement `internal/run/tui/app/instance/instance.go` table view with status indicators, columns, and navigation shortcuts
- [x] Task: Implement Instance Dashboard & Describe View (TDD) [0d11909]
    - [x] Write unit tests for Instance dashboard tabs and describe rendering in `internal/run/tui/app/instance/dashboard_test.go`
    - [x] Implement `internal/run/tui/app/instance/dashboard.go` (Overview, Containers, Conditions, YAML/JSON view)
- [x] Task: Implement Instance Lifecycle Actions & Modals (TDD) [043b3ec]
    - [x] Write unit tests for Start, Stop, and Delete action handlers/modals in `internal/run/tui/app/instance/action_test.go`
    - [x] Implement async Start, Stop, and Delete actions with spinners and confirmation modal in `internal/run/tui/app/instance/action.go`
- [x] Task: Conductor - User Manual Verification 'TUI Components for Instances (List, Dashboard, Lifecycle & Modals)' (Protocol in workflow.md)

## Phase 3: Application Integration & Documentation [checkpoint: 04dbd6f]
- [x] Task: Wire Instances Navigation and Shortcuts in Main App (TDD) [2fe8511]
    - [x] Write unit tests for `Ctrl+I` navigation, shortcut capture, and console URL generation in `internal/run/tui/app/app_test.go`
    - [x] Update `internal/run/tui/app/app.go` with `Ctrl+I` shortcut, console URL mapping, preload integration, and footer hints
- [x] Task: Update Centralized Help and Documentation (TDD/Docs) [e7c6825]
    - [x] Update help modal in `internal/run/tui/app/help/help.go` and tests in `internal/run/tui/app/help/help_test.go`
    - [x] Update `README.md` with the new Instances resource documentation and keyboard shortcuts
- [x] Task: Conductor - User Manual Verification 'Application Integration & Documentation' (Protocol in workflow.md)
