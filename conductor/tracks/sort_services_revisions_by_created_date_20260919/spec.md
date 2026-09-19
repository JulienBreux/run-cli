# Specification: Sort Services Revisions by Created Date

## Overview
Sort Cloud Run service revisions by creation date (`CreateTime`) in descending order (most recent first) within the API client layer (`internal/run/api/service/revision/revision.go`). This ensures all views and components consuming service revisions (such as the Service Dashboard revisions table and Traffic Split modal) display the latest revisions at the top.

This feature addresses and resolves GitHub issue #59.

## Functional Requirements
1. **API Layer Sorting**:
   - In `internal/run/api/service/revision/revision.go`, sort the `revisions` slice descending by `CreateTime` before returning from `List(...)`.
   - More recent revisions (later `CreateTime`) appear first.
   - For revisions with identical `CreateTime`, fall back to secondary sorting by revision `Name` descending.
   - Revisions with zero or missing `CreateTime` are placed at the very end of the list.
2. **Dashboard & TUI Integration**:
   - Maintain the revisions table header as `DEPLOYED` in `internal/run/tui/app/service/revision/revision.go`.
   - Ensure the Service Dashboard and Traffic Split components display the sorted revisions consistently.
3. **Documentation**:
   - Update `README.md` to document that service revisions are displayed in descending chronological order (most recent first).

## Non-Functional Requirements
- **Performance**: In-memory slice sorting overhead is minimal (typically <= 100 revisions per service).
- **Test Coverage**: Add comprehensive unit tests in `internal/run/api/service/revision/revision_test.go` covering descending sort, tie-breaking by Name, and zero timestamp handling. Maintain >80% test coverage.
- **Code Quality**: Passes `make lint` and compiles via `make build`.

## Acceptance Criteria
- [ ] `revision.List(...)` returns revisions sorted by `CreateTime` descending (most recent first).
- [ ] Ties in `CreateTime` are sorted by `Name` descending.
- [ ] Revisions with zero/unset `CreateTime` appear at the end of the list.
- [ ] Service Dashboard revisions table and traffic split modal reflect the sorted revisions.
- [ ] Resolves GitHub issue #59.
- [ ] `README.md` is updated to describe the revision ordering.
- [ ] Unit tests pass (`make test`).
- [ ] Linter passes (`make lint`).
- [ ] Binary builds successfully (`make build`).

## Out of Scope
- Interactive column sorting / re-sorting via keybindings in the TUI table.
- Modifying sorting behavior of non-revision resources (e.g., Jobs, Services, Instances).
