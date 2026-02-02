# Implementation Plan - Service Revision View

## Phase 1: Models and API Client
- [x] Task: Update Revision models and client to ensure all necessary fields for traffic and history are captured. 7ecda3e
    - [x] Write unit tests for Revision model parsing and client data retrieval.
    - [x] Implement/Update Revision model in internal/run/model/service/revision.
    - [x] Update API client in internal/run/api/service/revision.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Models and API Client' (Protocol in workflow.md)

## Phase 2: TUI Components
- [ ] Task: Create a dedicated Revision list component with traffic allocation visualization.
    - [ ] Write unit tests for Revision list component rendering logic.
    - [ ] Implement the Revision list table in internal/run/tui/app/service/revision.
- [ ] Task: Create a Revision detail view showing configuration and status.
    - [ ] Write unit tests for Revision detail component logic.
    - [ ] Implement the detail view component.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: TUI Components' (Protocol in workflow.md)

## Phase 3: Integration and Navigation
- [ ] Task: Integrate the new Revision views into the Service Dashboard.
    - [ ] Update Service Dashboard layout to include a Revisions tab/view.
    - [ ] Implement navigation between Service list, Dashboard, and Revision details.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Integration and Navigation' (Protocol in workflow.md)
