# Specification: Dynamic Terminal CLI Title

## Overview
Dynamically update the terminal window title in Run CLI to reflect the currently active context across both the interactive TUI application and CLI subcommands. Titles follow a consistent, hierarchical pattern starting with `"Run"`:
- Non-TUI Subcommands: `"Run | <Command>"` (e.g., `"Run | Version"` for the version command, resetting on exit).
- Initial TUI Startup: `"Run | Loading..."`
- Resource Lists: `"Run | <Resource>"` (e.g., `"Run | Services"`, `"Run | Jobs"`, `"Run | Instances"`, `"Run | Worker Pools"`, `"Run | Domain Mappings"`).
- Resource Dashboards/Details: `"Run | <Resource> | <Item Name>"` (e.g., `"Run | Services | frontend"`, `"Run | Jobs | backup-job"`, `"Run | Instances | instance-1"`).
- Modals & Overlays: Appends the modal context (e.g., `"Run | Services | Scale"`, `"Run | Instances | Authentication"`, `"Run | Help"`, `"Run | Switch Project"`, `"Run | Switch Region"`), restoring the underlying page's title upon closing.

## Functional Requirements
1. **Centralized Title Utility**:
   - Provide a title formatting and escape-sequence helper (e.g., in a dedicated package `internal/run/term` or `pkg/term`) supporting `FormatTitle(parts ...string) string` and ANSI terminal OSC escape sequences (`\033]0;%s\007`).
2. **Non-TUI CLI Subcommands**:
   - Update commands (such as `version`) to set the terminal title during execution (e.g., `"Run | Version"`) and restore/reset upon command termination.
3. **Interactive TUI Startup**:
   - Set the terminal title to `"Run | Loading..."` when the TUI application initiates the preloading screen.
4. **List View Navigation**:
   - Update title on switching to any resource list view (`Run | Services`, `Run | Jobs`, `Run | Instances`, `Run | Worker Pools`, `Run | Domain Mappings`).
5. **Detail Dashboard Navigation**:
   - Update title when opening resource details (`Run | <Resource> | <Item Name>`).
6. **Modals & Overlays Stack**:
   - Maintain title context when opening modals (Help modal, Scale modal, Authentication modal, Project switcher, Region switcher), pushing the modal title (e.g., `"Run | Help"`, `"Run | Services | Scale"`).
   - Restore the previous title when the modal is closed or dismissed.

## Non-Functional Requirements
- **Consistency**: Strict adherence to the `"Run | <Part 1> | <Part 2>"` pattern.
- **Terminal Compatibility**: Use standard OSC 0 sequences supported by xterm, iTerm2, Terminal.app, and modern terminal emulators.
- **Test Coverage**: Comprehensive unit tests covering title formatting, subcommand execution, and TUI title transitions. Maintain >80% coverage (aiming for 100%).
- **Code Quality**: Ensure all code passes `make test`, `make lint`, and `make build`.

## Acceptance Criteria
- [ ] Running `run version` sets terminal title to `"Run | Version"` and resets on exit.
- [ ] Starting the TUI displays `"Run | Loading..."` during initialization.
- [ ] Navigating between resource lists updates the terminal title accordingly (`"Run | Services"`, etc.).
- [ ] Navigating to detail dashboards updates the title to `"Run | <Resource> | <Item Name>"`.
- [ ] Modals append their name and restore the previous title upon closing.
- [ ] All unit tests pass (`make test`).
- [ ] Linter passes with 0 issues (`make lint`).
- [ ] Binary compiles cleanly (`make build`).
