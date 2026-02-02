# Implementation Plan - Help Page and Shortcut Reference

## Phase 1: Help Modal Component
- [x] Task: Create the Help Modal TUI component. 56263 (tests), 56359 (tests)
    - [x] Write unit tests for shortcut list data structure and category filtering if any.
    - [x] Implement the `HelpModal` in `internal/run/tui/app/help/help.go` using a `tview.Table` or `tview.TextView` with dynamic content.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Help Modal Component' (Protocol in workflow.md)

## Phase 2: Shortcut Registry and Mapping
- [x] Task: Define a centralized registry of shortcuts.
    - [x] Consolidate shortcut definitions (currently scattered in `app.go`, `service.go`, etc.) into a consistent data structure.
    - [x] Ensure all existing shortcuts (Global, Service, Job, Worker, Domain) are captured with their descriptions.
- [x] Task: Conductor - User Manual Verification 'Phase 2: Shortcut Registry and Mapping' (Protocol in workflow.md)

## Phase 3: Integration and Trigger
- [x] Task: Integrate the Help Modal into the main application flow.
    - [x] Add `?` key handling to the global input capture in `internal/run/tui/app/app.go`.
    - [x] Implement `openHelpModal` in `internal/run/tui/app/modal.go`.
    - [x] Update all footer shortcut hints to include `<?> Help`.
- [x] Task: Conductor - User Manual Verification 'Phase 3: Integration and Trigger' (Protocol in workflow.md)
