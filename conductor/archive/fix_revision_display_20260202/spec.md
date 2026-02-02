# Specification - Fix Stale Revision Display

## Overview
Address an issue where the Service Dashboard's revision display (list and details) retains data from previously viewed services or selections, leading to a confusing user experience.

## Functional Requirements
- **Clear Data on Exit:** When navigating away from the Service Dashboard (e.g., via `Esc`), all locally stored revision data for that service must be cleared.
- **Clear Data on New Selection:** When a new service is selected from the main list, any existing revision data in the dashboard components must be wiped before fetching the new data.
- **Visual State Management:**
    - The Revision list must be empty during transition or after exit.
    - The Revision details panel must be blank or show a loading indicator during transition.

## Acceptance Criteria
- [ ] Navigating from Service A to the List and then to Service B does not briefly show Service A's revisions.
- [ ] Pressing `Esc` to quit the Dashboard clears the `dashboardRevisions` and component text.
- [ ] Selecting a revision in Service A and then opening Service B results in an empty Details panel for Service B until a revision is selected.

## Out of Scope
- **API Enhancements:** Improving data fetching performance or modifying the actual Revision API client logic.
- **Other Dashboard Tabs:** Modifications to the Networking, Security, or Observability tabs (this track is strictly for Revision display).
- **Global State Management:** Large-scale refactoring of the application's global state or navigation framework.
- **Feature Additions:** Adding new revision-related features like tagging or traffic splitting (this is strictly a display fix).
