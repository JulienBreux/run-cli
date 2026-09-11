# Implementation Plan: Cloud Run Instances Authentication Management & Local Proxy Support

## Phase 1: API & Data Model Support (TDD) [checkpoint: 1365dce]
- [x] Task: Extend Instance Model with Proxy Status 9ce2ad0
    - [x] Write unit test for Instance model proxy field and status in `internal/run/model/instance/instance_test.go`
    - [x] Add `Proxy *model_service.ProxyStatus` to `Instance` struct in `internal/run/model/instance/instance.go`
    - [x] Run tests and verify coverage
- [x] Task: Implement `UpdateAuthentication` API for Instances 8030329
    - [x] Extend `InstancesClientWrapper` interface and mock in `internal/run/api/instance/client.go`
    - [x] Write unit tests for `UpdateAuthentication` in `internal/run/api/instance/instance_test.go`
    - [x] Implement `UpdateAuthentication` in `internal/run/api/instance/instance.go`
    - [x] Run tests and verify coverage
- [x] Task: Conductor - User Manual Verification 'API & Data Model Support' (Protocol in workflow.md)

## Phase 2: TUI List View & Proxy Management (TDD)
- [x] Task: Update Instances Table Headers and Cells for PROXY and AUTH 6072d28
    - [x] Write unit tests for table rendering with proxy and auth statuses in `internal/run/tui/app/instance/instance_test.go`
    - [x] Update `listHeaders` and `listExpansions` with `PROXY` and `AUTH` columns
    - [x] Render proxy status indicator (`[green]P[white]`) and auth status indicator (`[red]Yes[white]` / `[green]No[white]`)
    - [x] Run tests and verify coverage
- [~] Task: Implement Proxy Toggle and Open URL Actions for Instances
    - [ ] Write unit tests for `toggleProxy` and `OpenURL` logic in `internal/run/tui/app/instance/instance_test.go`
    - [ ] Integrate `proxy.Manager` in `internal/run/tui/app/instance/instance.go` to start/stop proxy on instance primary URL
    - [ ] Implement `GetSelectedInstanceURL()` and open browser URL logic (`o` key)
    - [ ] Run tests and verify coverage
- [ ] Task: Register Shortcuts and Update Context Footer
    - [ ] Register `a`, `p`, and `o` in `internal/run/tui/app/shortcut/shortcut.go` under `CategoryInstanceList`
    - [ ] Update dynamic shortcut bar in `instance.Shortcuts()` to show active proxy port when enabled
    - [ ] Run tests in `shortcut_test.go` and `instance_test.go`
- [ ] Task: Conductor - User Manual Verification 'TUI List View & Proxy Management' (Protocol in workflow.md)

## Phase 3: Authentication Modal & App Integration (TDD)
- [ ] Task: Create Instance Authentication Modal
    - [ ] Write unit tests for Instance Auth modal in `internal/run/tui/app/instance/auth/auth_test.go`
    - [ ] Implement modal in `internal/run/tui/app/instance/auth/auth.go` with dropdown and asynchronous save
    - [ ] Run tests and verify coverage
- [ ] Task: Wire Keyboard Events and Modal in Application Controller
    - [ ] Add `openInstanceAuthModal` in `internal/run/tui/app/modal.go`
    - [ ] Wire `a`, `p`, and `o` shortcuts in `internal/run/tui/app/app.go` for `instance.LIST_PAGE_ID`
    - [ ] Write unit tests in `internal/run/tui/app/app_test.go` and `modal_test.go`
    - [ ] Run tests and verify coverage
- [ ] Task: Conductor - User Manual Verification 'Authentication Modal & App Integration' (Protocol in workflow.md)

## Phase 4: Verification & Documentation
- [ ] Task: Document New Shortcuts in README.md
    - [ ] Update `README.md` with Instance proxy (`p`), auth (`a`), and open (`o`) shortcuts and table columns
- [ ] Task: Quality Gates & Linting
    - [ ] Run `make test` to verify complete test suite and coverage
    - [ ] Run `make lint` to ensure zero lint errors
    - [ ] Run `make build` to verify clean compilation
- [ ] Task: Conductor - User Manual Verification 'Verification & Documentation' (Protocol in workflow.md)
