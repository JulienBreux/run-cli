# Specification: Clean Maintainer and Community Health Files

## Overview
Address issue [#51](https://github.com/JulienBreux/run-cli/issues/51) and align the repository with GitHub's community health standards. Separate maintainer rosters from contribution instructions by establishing a dedicated `MAINTAINERS.md`, rewriting `CONTRIBUTING.md` into an actionable guide for contributors, and adding default community health files (`CODE_OF_CONDUCT.md`, `SECURITY.md`, `SUPPORT.md`, and `.github/pull_request_template.md`).

## Functional Requirements
1. **Maintainers Roster (`MAINTAINERS.md`)**:
   - Extract the maintainers and contributors table currently residing in `CONTRIBUTING.md` into a dedicated `MAINTAINERS.md` at repository root.
   - List Top-level maintainers/owners (Julien Breux) and contributors.
   - Clearly describe roles, ownership, and how to become a maintainer, linking directly to `CONTRIBUTING.md`.

2. **Contribution Guidelines (`CONTRIBUTING.md`)**:
   - Replace `CONTRIBUTING.md` with a comprehensive guide for human contributors.
   - Include sections on:
     - Welcome & Code of Conduct reference.
     - Getting Started & Prerequisites (Go 1.27+, Docker, Make, gcloud).
     - Development commands (`make build`, `make run`, `make test`, `make lint`).
     - Coding and testing conventions (TDD, test coverage, Apache 2.0 license headers on source files, Conventional Commits).
     - Pull request submission workflow and branch naming.
     - References to `MAINTAINERS.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, and `SUPPORT.md`.

3. **Community Health & Support Documentation**:
   - **`CODE_OF_CONDUCT.md`**: Contributor Covenant v2.1 standard with reporting contact.
   - **`SECURITY.md`**: Security vulnerability disclosure process and supported versions.
   - **`SUPPORT.md`**: Guidance for seeking help, reporting bugs, and utilizing GitHub Discussions / Issues.
   - **`.github/pull_request_template.md`**: Structured PR checklist (description, linked issues, type of change, testing verification).

4. **Repository Documentation Alignment (`README.md`)**:
   - Update the "Contributing" section of `README.md` to link to `CONTRIBUTING.md`, `MAINTAINERS.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, and `SUPPORT.md`.

## Non-Functional Requirements
- **Consistency**: Build and test instructions must match the existing `Makefile` targets (`make build`, `make test`, `make lint`, `make run`).
- **Link Integrity**: All relative links across documentation files must be valid and resolve properly.
- **GitHub Compliance**: Adheres to GitHub's community health file naming and placement standards.

## Acceptance Criteria
- [ ] `MAINTAINERS.md` contains the maintainers and contributors table and links to `CONTRIBUTING.md`.
- [ ] `CONTRIBUTING.md` contains complete development and contribution workflows without recursive self-links.
- [ ] `CODE_OF_CONDUCT.md`, `SECURITY.md`, and `SUPPORT.md` exist at repository root with accurate project details.
- [ ] `.github/pull_request_template.md` exists and enforces quality checklists for pull requests.
- [ ] `README.md` has updated links to community health documents.
- [ ] All cross-references and links are verified.
- [ ] Issue #51 is resolved.

## Out of Scope
- Code changes to Go packages or TUI components.
- Automated bot configurations (e.g., CLA bots, stale bots).
