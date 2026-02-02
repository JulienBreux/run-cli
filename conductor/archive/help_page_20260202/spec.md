# Specification - Help Page and Shortcut Reference

## Overview
Implement a comprehensive help page accessible via the `?` keyboard shortcut. This page will provide users with a clear reference of all global and contextual keyboard shortcuts available within the Run CLI TUI, organized by functional category.

## Functional Requirements
- **Shortcut Trigger:** Pressing the `?` key from any active view must open the Help Modal.
- **Help Modal:**
    - Displayed as a centralized modal overlay.
    - Dismissible via `Esc`, `Enter`, or pressing `?` again.
- **Categorized Shortcut List:**
    - Shortcuts must be grouped by context:
        - **Global:** (e.g., Switching between views, opening project/region modals).
        - **Service List:** (e.g., Describe, Logs, Traffic Split).
        - **Service Dashboard:** (e.g., Tab navigation, Traffic Split).
        - **Job List:** (e.g., Execution).
        - **Worker List:** (e.g., Scaling).
- **Entry Format:** Each shortcut entry must display the key combination followed by a clear description of the action.

## Non-Functional Requirements
- **Visual Consistency:** Use existing TUI modal components and consistent typography.
- **Performance:** The modal must open near-instantaneously without any significant lag.

## Acceptance Criteria
- [ ] Pressing `?` from any screen opens the Help Modal.
- [ ] The Help Modal displays a complete list of categorized shortcuts.
- [ ] The descriptions match the actual behavior of the keys.
- [ ] The modal can be closed easily to return to the previous context.

## Out of Scope
- Detailed user documentation or tutorials (strictly a shortcut reference).
- Interactive tutorials or onboarding flows.
- Modifying or adding new shortcuts (focus is on documentation of existing ones).
