# Implementation Plan - Service Revision View

## Phase 1: Models and API Client [checkpoint: 4560608]
- [x] Task: Update Revision models and client to ensure all necessary fields for traffic and history are captured. 7ecda3e
    - [x] Write unit tests for Revision model parsing and client data retrieval.
    - [x] Implement/Update Revision model in internal/run/model/service/revision.
    - [x] Update API client in internal/run/api/service/revision.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Models and API Client' (Protocol in workflow.md) 4560608

## Phase 2: TUI Components [checkpoint: e863d9d]
- [x] Task: Create a dedicated Revision list component with traffic allocation visualization. 98ae517
    - [x] Write unit tests for Revision list component rendering logic.
    - [x] Implement the Revision list table in internal/run/tui/app/service/revision.
- [x] Task: Create a Revision detail view showing configuration and status. a144f5d
    - [x] Write unit tests for Revision detail component logic.
    - [x] Implement the detail view component.
- [x] Task: Conductor - User Manual Verification 'Phase 2: TUI Components' (Protocol in workflow.md) e863d9d

## Phase 3: Integration and Navigation [checkpoint: 5e92f29]
- [x] Task: Integrate the new Revision views into the Service Dashboard. bb96a70
    - [x] Update Service Dashboard layout to include a Revisions tab/view.
    - [x] Implement navigation between Service list, Dashboard, and Revision details.
- [x] Task: Conductor - User Manual Verification 'Phase 3: Integration and Navigation' (Protocol in workflow.md) 5e92f29
