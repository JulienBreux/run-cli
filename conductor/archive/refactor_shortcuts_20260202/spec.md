# Specification - Refactor Shortcut Display

## Overview
Refactor the Text-based User Interface (TUI) to utilize the centralized `shortcut` package for all shortcut displays, including footer hints across all functional views (Services, Jobs, Worker Pools, and Domain Mappings). This will eliminate hardcoded strings across the codebase, ensuring architectural consistency and a single source of truth for keyboard commands.

## Functional Requirements
- **Centralized Registry:** Ensure all contextual shortcuts for Services, Jobs, Worker Pools, and Domain Mappings are present in `internal/run/tui/app/shortcut/shortcut.go`.
- **Dynamic Formatting Helper:** Implement a function in the `shortcut` package to generate standardized footer strings (e.g., `[dodgerblue]<key> [white]Desc`) for a given category.
- **Dynamic Shortcut Support:** Implement a mechanism to handle shortcuts with state-dependent descriptions (e.g., Proxy or Authentication toggles).
- **View Integration:** Update all functional views to derive their shortcut hints from the `shortcut` package.

## Acceptance Criteria
- [ ] No hardcoded shortcut hint strings exist in `internal/run/tui/app/...`.
- [ ] Footer shortcuts for all views are generated using the `shortcut` package.
- [ ] Dynamic shortcuts (e.g., Proxy, Authentication status) continue to work and display correctly.
- [ ] The TUI behavior remains identical to the current implementation but with improved internal structure.

## Out of Scope
- Adding new shortcuts or changing existing key bindings.
- Large-scale visual changes to the footer or header layout.
