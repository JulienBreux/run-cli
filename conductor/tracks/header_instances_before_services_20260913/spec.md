# Specification: Place Instances before Services in Header

## Overview
In the Run CLI terminal user interface, reorder the global shortcuts so that "Instances" (`<ctrl+n>`) appears before "Services" (`<ctrl+s>`). By updating the canonical global shortcut registry, this reordering is reflected across both the Header component and the Help modal (`?`).

## Functional Requirements
- **FR-1:** In `internal/run/tui/app/shortcut/shortcut.go`, reorder the `Registry` entries under `CategoryGlobal` such that `Instances` (`ctrl+n`) is placed immediately before `Services` (`ctrl+s`).
- **FR-2:** The Header component (`internal/run/tui/component/header/header.go`) must display `<ctrl+n> Instances` above `<ctrl+s> Services` in the global shortcuts column.
- **FR-3:** The Help modal (`?`) must display `<ctrl+n> Instances` before `<ctrl+s> Services` in the Global category list.
- **FR-4:** Preserve the existing label "Instances" (plural) to stay consistent with other resource categories ("Services", "Jobs", "Worker Pools").

## Non-Functional Requirements
- **NFR-1:** All existing unit tests in `shortcut_test.go`, `header_test.go`, and related packages must continue to pass.
- **NFR-2:** Code coverage must remain above 80% with no lint regressions (`make test`, `make lint`).

## Acceptance Criteria
- [ ] In the Header component's shortcut list, `<ctrl+n> Instances` appears before `<ctrl+s> Services`.
- [ ] In the Help modal (`?`), `<ctrl+n> Instances` appears before `<ctrl+s> Services`.
- [ ] Unit tests verify the order of shortcuts in `shortcut.Registry` and ensure header rendering behavior is preserved.
- [ ] `make test` and `make lint` pass without errors.

## Out of Scope
- Changing the initial landing screen on application launch (remains Services).
- Modifying keybindings or behavior for `<ctrl+n>` and `<ctrl+s>`.
- Renaming "Instances" to singular "Instance".
