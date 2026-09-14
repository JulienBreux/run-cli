# Specification: Resolve Version Information for go install Builds and Check for Updates

## Overview
When `run-cli` is installed via `go install github.com/JulienBreux/run-cli/cmd/run@<version>`, the binary does not receive `-ldflags` compiler flags typically injected by GoReleaser or the `Makefile`. Currently, this results in version output defaulting to `dev` with commit `n/a` and date `n/a`. By integrating Go's `runtime/debug.ReadBuildInfo()`, the application will automatically extract and populate the module version, VCS commit hash, and build timestamp when compile-time ldflags are unset. Additionally, an asynchronous non-blocking remote release check will query GitHub releases in the background to notify the user if an update is available.

## Functional Requirements
- **FR-1:** If `Version == "dev"` (compile-time ldflags not provided), query `runtime/debug.ReadBuildInfo()` to detect version and VCS build metadata.
- **FR-2:** If `info.Main.Version` is present and valid, assign it to `Version`. If `info.Main.Version` is `(devel)` or empty, retain `dev`.
- **FR-3:** If `Commit == "n/a"`, attempt to read `vcs.revision` from build info settings and set `Commit`. If `vcs.modified` is true, append `-dirty`.
- **FR-4:** If `RawDate == "n/a"`, attempt to read `vcs.time` from build info settings and set `RawDate`.
- **FR-5:** Preserve precedence for explicit compile-time flags passed via `-ldflags -X github.com/JulienBreux/run-cli/pkg/version.Version=...`.
- **FR-6:** Provide an update checker in `pkg/version` that queries GitHub's latest release API (`https://api.github.com/repos/JulienBreux/run-cli/releases/latest`) asynchronously with a short timeout (e.g., 2s) so it never hangs offline or in slow environments.
- **FR-7:** In the TUI Header (`internal/run/tui/component/header`), if an update is detected, display a subtle notification beside the version (e.g., `v0.2.0 (update: v0.3.0)`).
- **FR-8:** In `run version` CLI command, if update information is retrieved or requested, report whether a newer version is available.

## Non-Functional Requirements
- **NFR-1:** Asynchronous and non-blocking: network checks must run in a background goroutine and never delay or degrade CLI startup or UI rendering.
- **NFR-2:** Robust failure handling: network timeouts, offline mode, rate limits, or API errors must fail silently without errors or crashes.
- **NFR-3:** Maintain code coverage above 80% for `pkg/version` and affected packages.
- **NFR-4:** All existing unit tests and linters (`make test`, `make lint`) must pass cleanly.

## Acceptance Criteria
- [ ] Running `run version` on binaries built via `go install` reflects the correct module version (e.g. `v0.2.0`), commit hash, and build date.
- [ ] Binaries built with explicit `-ldflags` continue to honor injected values with highest priority.
- [ ] TUI Header updates dynamically in the background to show update availability when a newer release is published.
- [ ] Unit tests verify build info extraction, fallback logic, and release comparison using mocked responses.
- [ ] Quality gates (`make test`, `make lint`, `make build`) pass with zero regressions.

## Out of Scope
- Automatic self-updating or replacing the binary.
- Blocking the UI waiting for network responses.
