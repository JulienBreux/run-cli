# Specification: Exclude Conductor Commits from GoReleaser Changelog & Update Release Workflow

## Overview
Update the GoReleaser changelog configuration in `.goreleaser.yaml` to exclude all Conductor-related commits (`chore(conductor): ...`, `conductor(plan): ...`, etc.) and support Conventional Commit scopes for `docs` and `test` exclusions. Additionally, enhance `.github/workflows/release.yml` with a manual `workflow_dispatch` trigger while preserving existing release credentials.

## Functional Requirements
1. **Exclude Conductor Commits**:
   - Update `changelog.filters.exclude` in `.goreleaser.yaml` to exclude any commit matching "conductor" (case-insensitive `(?i)conductor`).
2. **Support Scoped Docs & Tests**:
   - Update exclusion patterns for documentation (`^docs(\(.*\))?:`) and tests (`^test(\(.*\))?:`) so that scoped commits like `docs(maintainers): ...` or `test(unit): ...` are filtered out.
3. **Workflow Dispatch Trigger**:
   - Add `workflow_dispatch:` trigger to `.github/workflows/release.yml` to enable manual triggering of the release workflow alongside tag pushes and pull requests.
   - Maintain existing credentials (`secrets.GH_TOKEN` for `GITHUB_TOKEN` and `HOMEBREW_GITHUB_API_TOKEN`, `secrets.CODECOV_TOKEN`).
4. **Configuration Validation**:
   - Validate `.goreleaser.yaml` using `goreleaser check`.

## Non-Functional Requirements
- **Changelog Clarity**: Ensure generated changelogs are clean and focused on user-facing changes.
- **Workflow Reliability**: Verify YAML syntax and GitHub Actions workflow validity.

## Acceptance Criteria
- [ ] `.goreleaser.yaml` filters exclude `(?i)conductor` commits.
- [ ] `.goreleaser.yaml` filters exclude scoped `docs(...)` and `test(...)` commits.
- [ ] `goreleaser check` succeeds with 0 errors.
- [ ] `.github/workflows/release.yml` includes `workflow_dispatch` trigger.
- [ ] All existing project tests and linters pass (`make test`, `make lint`, `make build`).

## Out of Scope
- Modifying release artifact formats, targets, or binary packaging settings.
