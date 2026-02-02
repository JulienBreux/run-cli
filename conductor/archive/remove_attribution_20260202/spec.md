# Specification - Remove Author Attribution from Header

## Overview
Remove the "♥ Julien Breux" author attribution specifically from the global TUI header. This change aims to provide a more neutral, professional appearance during regular tool usage and optimize the header layout by reclaiming horizontal space.

## Functional Requirements
- **Header Attribution Removal:** Permanently delete the "♥ Julien Breux" text from the global header component.
- **Layout Adjustment:** Ensure the header layout remains balanced and aesthetically pleasing after the removal. The space previously occupied by the mention should either remain empty or be absorbed by neighboring components (Logo or Shortcuts).

## Non-Functional Requirements
- **Visual Integrity:** Maintain the overall design language and alignment of the header.

## Acceptance Criteria
- [ ] The "♥ Julien Breux" text is no longer visible in the global header across all views.
- [ ] The header correctly renders without errors or misalignments.
- [ ] The attribution remains visible in the **Credits** modal.
- [ ] The attribution remains visible in the **Startup Logo/Loading** screen.

## Out of Scope
- Modifying the "Credits" or "About" modal.
- Modifying the startup logo or initial loading animation.
- Adding new header features or information.
