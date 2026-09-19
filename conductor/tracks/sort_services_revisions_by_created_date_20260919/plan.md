# Implementation Plan - Sort Services Revisions by Created Date

## Phase 1: API Layer Revision Sorting (TDD)
- [ ] Task: Write Failing Unit Tests for Revision Sorting (Red Phase)
    - [ ] Add unit tests in `internal/run/api/service/revision/revision_test.go` asserting revisions returned by `List(...)` are ordered by `CreateTime` descending (most recent first)
    - [ ] Add tests verifying tie-breaker sorting by revision `Name` descending when timestamps match
    - [ ] Add tests verifying revisions with zero/unset `CreateTime` are placed at the end of the list
- [ ] Task: Implement Revision Sorting in API Client (Green Phase)
    - [ ] Implement descending sort logic in `internal/run/api/service/revision/revision.go`
    - [ ] Verify all tests pass with `go test` / `make test` and module coverage remains >80%
- [ ] Task: Conductor - User Manual Verification 'Phase 1: API Layer Revision Sorting (TDD)' (Protocol in workflow.md)

## Phase 2: Documentation & Full Verification
- [ ] Task: Documentation & Issue Resolution
    - [ ] Update `README.md` to document that service revisions are displayed in descending chronological order (most recent first)
    - [ ] Document resolution of GitHub issue #59
- [ ] Task: Full Suite Verification & Build
    - [ ] Run `make test` to ensure all unit tests pass
    - [ ] Run `make lint` to verify zero linting or static analysis issues
    - [ ] Run `make build` to verify clean compilation of the `./bin/run` binary
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Documentation & Full Verification' (Protocol in workflow.md)
