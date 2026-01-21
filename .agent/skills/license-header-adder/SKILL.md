---
name: license-header-adder
description: Adds the standard open-source license header to new source files. Use involves creating new code files that require copyright attribution.
---

# License Header Adder Skill

This skill ensures that all new source files have the correct copyright header.

## Instructions

1. **Read the Template**:
  First, read the content of the header template file located at `resources/HEADER_TEMPLATE.txt`.

2. **Prepend to File**:
  When creating a new file (e.g., `.go`), prepend the `target_file` content with the template content and change with the current year.

3. **Modify Comment Syntax**:
  - For C-style languages (Go), keep the `/* ... */` block as is.
  - For Shell or YAML, convert the block to use `#` comments.
