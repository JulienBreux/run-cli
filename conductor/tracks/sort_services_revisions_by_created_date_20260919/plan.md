# Implementation Plan - Sort Services Revisions by Created Date

## Phase 1: API Layer Revision Sorting (TDD) [checkpoint: 3ff78b7]
- [x] Task: Write Failing Unit Tests for Revision Sorting (Red Phase) 602ec7d
    - [x] Add unit tests in `internal/run/api/service/revision/revision_test.go` asserting revisions returned by `List(...)` are ordered by `CreateTime` descending (most recent first)
    - [x] Add tests verifying tie-breaker sorting by revision `Name` descending when timestamps match
    - [x] Add tests verifying revisions with zero/unset `CreateTime` are placed at the end of the list
- [x] Task: Implement Revision Sorting in API Client (Green Phase) e4755dd
    - [x] Implement descending sort logic in `internal/run/api/service/revision/revision.go`
    - [x] Verify all tests pass with `go test` / `make test` and module coverage remains >80%
- [x] Task: Conductor - User Manual Verification 'Phase 1: API Layer Revision Sorting (TDD)' (Protocol in workflow.md) 3ff78b7

## Phase 2: Documentation & Full Verification
- [x] Task: Documentation & Issue Resolution f173b88
    - [x] Update `README.md` to document that service revisions are displayed in descending chronological order (most recent first)
    - [x] Document resolution of GitHub issue #59
- [x] Task: Full Suite Verification & Build
    - [x] Run `make test` to ensure all unit tests pass
    - [x] Run `make lint` to verify zero linting or static analysis issues
    - [x] Run `make build` to verify clean compilation of the `./bin/run` binary
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Documentation & Full Verification' (Protocol in workflow.md)
