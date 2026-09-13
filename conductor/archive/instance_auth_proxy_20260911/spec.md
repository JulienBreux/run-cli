# Specification: Cloud Run Instances Authentication Management & Local Proxy Support

## Overview
Add authentication management and local proxy features to Cloud Run Instances in `run-cli`, matching the existing Service capabilities:
1. **IAM Authentication Management (`a` shortcut)**: Toggle between requiring authentication and allowing unauthenticated invocations.
2. **Local Authenticated Proxy (`p` shortcut)**: Run a local HTTP reverse proxy that injects Google ID token credentials when communicating with the instance URL.
3. **Open URL (`o` shortcut)**: Open the instance URL (or active local proxy URL) in the default browser.

## Functional Requirements
1. **API Client (`internal/run/api/instance`)**:
   - Implement `UpdateAuthentication(ctx context.Context, project, region, instanceName string, allowUnauthenticated bool) (*model.Instance, error)` to update `InvokerIamDisabled`.
   - Update `InstancesClientWrapper` and mock implementations to support testing the authentication update.

2. **Data Model (`internal/run/model/instance`)**:
   - Add `Proxy *model_service.ProxyStatus` to `model.Instance` to maintain local proxy status (Enabled, Port, URL).

3. **TUI List Table (`internal/run/tui/app/instance`)**:
   - Add `PROXY` and `AUTH` columns to the Instances list table.
   - Display `[green]P[white]` when local proxy is running.
   - Display `[red]Yes[white]` when authentication is required and `[green]No[white]` when unauthenticated access is permitted (`InvokerIamDisabled` is true).

4. **Interactive Actions & Shortcuts**:
   - **Authentication Modal (`a`)**: Open a modal dialog to toggle between "Require authentication" and "Allow unauthenticated invocations" with animated spinner and error handling.
   - **Toggle Proxy (`p`)**: Start or stop a local reverse proxy targeting the instance's primary URL (`instance.URLs[0]`) using `proxy.Manager`, with live port display in footer shortcut hints.
   - **Open URL (`o`)**: Open the active instance URL (local proxy URL if enabled, otherwise cloud URL) in the system browser.

5. **Shortcuts & Documentation**:
   - Register `a`, `p`, and `o` in `shortcut.Registry` under `CategoryInstanceList`.
   - Wire input handlers in `internal/run/tui/app/app.go` and `internal/run/tui/app/instance/instance.go`.
   - Update `README.md` with the new shortcuts.

## Acceptance Criteria
- [ ] Table headers for Instances include `PROXY` and `AUTH`.
- [ ] Pressing `a` on an instance opens the Authentication modal and updates IAM settings via API.
- [ ] Pressing `p` starts/stops the local authenticated proxy for the instance, showing green `P` and port in footer.
- [ ] Pressing `o` opens the active URL or local proxy URL in the browser.
- [ ] Unit tests pass with >80% coverage on new/modified code; `make test`, `make lint`, and `make build` pass.
- [ ] `README.md` is updated.
