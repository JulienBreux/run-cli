# Implementation Plan: Increase Test Coverage to >80%

This plan outlines the steps to identify low-coverage areas and systematically increase the test coverage of the `run-cli` project to meet the >80% requirement.

## Phase 1: Assessment and Tooling Configuration [checkpoint: 15b3ebb]
Goal: Identify current coverage gaps and ensure reporting is accurate.

- [x] Task: Identify current coverage for all packages using `make test`. e252a58
- [x] Task: Document the packages and files with coverage below 80%. e252a58
- [x] Task: Conductor - User Manual Verification 'Assessment and Tooling Configuration' (Protocol in workflow.md) 15b3ebb

## Phase 2: Core Logic and Models (`internal/run/model`)
Goal: Increase coverage for the foundational data structures and business logic.

- [ ] Task: Write Tests: Add unit tests for `internal/run/model/common/...`.
- [ ] Task: Implement: Ensure all core model logic is verified.
- [ ] Task: Write Tests: Add unit tests for `internal/run/model/service/...`, `job/...`, and `workerpool/...`.
- [ ] Task: Implement: Enhance coverage for service, job, and workerpool models.
- [ ] Task: Conductor - User Manual Verification 'Core Logic and Models' (Protocol in workflow.md)

## Phase 3: API Clients (`internal/run/api`)
Goal: Ensure robust interactions with Google Cloud APIs through mocking.

- [ ] Task: Write Tests: Enhance test coverage for `internal/run/api/service`.
- [ ] Task: Implement: Add mocks for GCP Service API and verify error handling.
- [ ] Task: Write Tests: Enhance test coverage for `internal/run/api/job` and `internal/run/api/workerpool`.
- [ ] Task: Implement: Add mocks and fixtures for Job and Workerpool APIs.
- [ ] Task: Write Tests: Enhance test coverage for `internal/run/api/project` and `internal/run/api/region`.
- [ ] Task: Implement: Verify discovery logic for projects and regions.
- [ ] Task: Conductor - User Manual Verification 'API Clients' (Protocol in workflow.md)

## Phase 4: CLI Commands and TUI Components (`internal/run/command`, `internal/run/tui`)
Goal: Verify the user interface and command-line entry points.

- [ ] Task: Write Tests: Add unit tests for `internal/run/command` logic.
- [ ] Task: Implement: Ensure command flag parsing and basic execution flows are tested.
- [ ] Task: Write Tests: Increase coverage for TUI components in `internal/run/tui/component`.
- [ ] Task: Implement: Verify rendering and interaction logic for key UI elements.
- [ ] Task: Write Tests: Increase coverage for TUI application pages in `internal/run/tui/app`.
- [ ] Task: Implement: Test page transitions and data display.
- [ ] Task: Conductor - User Manual Verification 'CLI Commands and TUI Components' (Protocol in workflow.md)

## Phase 5: Final Verification and Reporting
Goal: Confirm the overall project coverage meets the goal.

- [ ] Task: Run final coverage report using `make test`.
- [ ] Task: Verify that the overall project coverage is >80%.
- [ ] Task: Update `README.md` or documentation if coverage requirements changed.
- [ ] Task: Conductor - User Manual Verification 'Final Verification and Reporting' (Protocol in workflow.md)
