# Specification - Update Traffic Split

## Overview
Implement the capability to update traffic split percentages between multiple revisions of a Cloud Run service directly within the Run CLI TUI. This enables controlled rollouts and A/B testing.

## Functional Requirements
- **Revision Selection:** Users can select multiple revisions from the existing Revision list table.
- **Traffic Update Modal:** An interactive modal triggered for the selected revisions where users can enter percentage values for each.
- **Real-time Validation:** The modal must ensure that the sum of all assigned percentages equals exactly 100% before allowing submission.
- **GCP API Integration:** Use the Google Cloud Run API to update the service's traffic configuration.
- **Error Handling:** Display descriptive error messages in a modal if the backend API rejects the update (e.g., due to permission issues or invalid revision states).

## Non-Functional Requirements
- **Responsive UI:** The modal should be responsive and centered within the TUI.
- **Visual Consistency:** Use existing TUI components (forms, inputs, buttons) and styling.

## Acceptance Criteria
- [ ] A user can select 2 or more revisions and trigger a "Traffic Split" action.
- [ ] The modal correctly lists the selected revisions and provides input fields for percentages.
- [ ] The "Save" button is disabled or triggers a local error if the sum is not 100%.
- [ ] Successfully saving updates the service's traffic split on GCP.
- [ ] The Revision list reflects the new traffic splits after the operation completes.

## Out of Scope
- Assigning or modifying revision tags (handled in a future track).
- Fast rollback to previous revision (handled in a future track).
