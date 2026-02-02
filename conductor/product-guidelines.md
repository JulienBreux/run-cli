# Product Guidelines

## Documentation & Communication
- **Technical & Precise Prose:** All documentation, tooltips, and user-facing messages must prioritize technical accuracy and precision. The language should be direct and consistent with cloud infrastructure terminology.

## Visual Identity & UI Design
- **Clarity & Information Density:** The TUI layout should be optimized to present the maximum amount of relevant data clearly. Avoid unnecessary whitespace or decorative elements that do not serve an informational purpose.
- **Professional Branding:** Maintain a clean and professional appearance by avoiding personal attributions in persistent UI elements like the header. Personal recognition is reserved for dedicated sections like the Credits modal.
- **High Contrast & Accessibility:** Use standard terminal colors to ensure high contrast and legibility across various terminal themes. Prioritize accessibility by avoiding color-only indicators where possible.

## User Interaction
- **Immediate Feedback:** Every asynchronous or long-running operation (e.g., fetching resources, scaling) must be accompanied by a visual indicator, such as a spinner or loader, to confirm the application is active.

## Performance & Responsiveness
- **Snappy Interface:** Prioritize near-instant UI updates. Transitions between pages and response to user input must feel immediate to maintain the high-performance feel of a local CLI tool.
