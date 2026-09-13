# Specification: Cloud Run Instances Resource Support

## Overview
Introduce support for the new Google Cloud Run "Instances" resource (`projects.locations.instances` in Cloud Run API v2) into `run-cli`. This brings full visibility and lifecycle management for individual Cloud Run instances alongside Services, Jobs, and Worker Pools.

## Functional Requirements
1. **API Client & Integration (`internal/run/api/instance`)**:
   - Integrate with `cloud.google.com/go/run/apiv2` using `InstancesClient` and `runpb.Instance`.
   - Implement `InstancesClientWrapper` with interface abstraction for mockable unit testing.
   - Support `ListInstances`, `GetInstance`, `StartInstance`, `StopInstance`, and `DeleteInstance` API operations.
   - Handle region filtering (`all` vs specific region) and project scoping.

2. **Data Models (`internal/run/model/instance`)**:
   - Model `Instance` struct encapsulating fields from `runpb.Instance`: ID, Name, Region, Project, CreateTime, UpdateTime, RestartPolicy, ContainerStatuses, Conditions, TerminalCondition, LogURI, URLs, Etag, Labels, ServiceAccount, etc.
   - Conversion functions between protobuf `runpb.Instance` and domain `model.Instance`.
   - Helper methods to evaluate overall instance health, status indicators, and display strings.

3. **TUI Views & Navigation (`internal/run/tui/app/instance`)**:
   - **List View**: Interactive table displaying instances with columns: Name, Region, Status, Restart Policy, Containers, Last Updated.
   - **Navigation**: Global shortcut `Ctrl+I` to navigate directly to the Instances page from anywhere in the TUI.
   - **Dashboard / Describe View**: Detailed tabbed/structured view showing Overview, Container Statuses & Resources, Conditions & Networking, and Raw YAML/JSON representation.
   - **Lifecycle Actions**:
     - Start instance (`s`) with spinner and status update.
     - Stop instance (`x`) with spinner and status update.
     - Delete instance (`d`) with confirmation modal dialog.
   - **Log Integration (`l`)**: Quick view of Cloud Logging logs associated with the instance.
   - **Cloud Console Integration (`Ctrl+Z`)**: Open the selected instance's console page in the default web browser.

4. **Application Wiring & Integration**:
   - Register `Ctrl+I` in global shortcuts and add to the footer help hints.
   - Update centralized help modal (`?`) with Instances shortcuts and actions.
   - Integrate into the application preloader for parallel data fetching alongside services and jobs.
   - Update `README.md` to document the Instances feature and keybindings.

## Non-Functional Requirements
- **TDD & High Coverage**: Unit tests written for all new packages (API, model, TUI logic), ensuring test coverage >80%.
- **Async & Non-blocking**: API calls and mutations execute asynchronously using existing `spinner` and `loader` components to ensure UI responsiveness.
- **Code Quality**: Passes `make lint` (`golangci-lint`) and `make test`. All new Go files contain the project's standard Apache 2.0 license header.

## Acceptance Criteria
- [ ] Pressing `Ctrl+I` opens the Instances list page populated with instances for the current project/region.
- [ ] Selecting an instance and pressing Enter opens the Instance dashboard view.
- [ ] User can trigger Start, Stop, and Delete actions with appropriate confirmation and visual progress.
- [ ] Pressing `l` views instance logs and `Ctrl+Z` opens the Cloud Console URL.
- [ ] Help modal (`?`) accurately lists all new shortcuts.
- [ ] `make test` passes with >80% coverage on new code; `make lint` and `make build` pass cleanly.
- [ ] `README.md` is updated.

## Out of Scope
- Standalone CLI command-line subcommands (e.g. `run instances ...`).
- Interactive instance creation wizard inside TUI (creation is managed externally via CI/CD, Terraform, or gcloud).
