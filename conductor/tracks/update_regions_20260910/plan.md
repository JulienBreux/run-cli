# Implementation Plan: Add Makefile Command to Update Cloud Run Regions

## Phase 1: Unit Testing Baseline & Validation
- [x] Task: Enhance Region Unit Tests (TDD) [b00c03b]
    - [x] Add unit tests verifying region format, non-empty list, and known standard regions in `internal/run/api/region/region_test.go`
    - [x] Verify test suite passes with `go test ./internal/run/api/region/...`
- [ ] Task: Conductor - User Manual Verification 'Unit Testing Baseline & Validation' (Protocol in workflow.md)

## Phase 2: Makefile Implementation & Region Update Command
- [ ] Task: Implement `regions-update` Target in Makefile
    - [ ] Add prerequisite validation ensuring `gcloud` is installed
    - [ ] Implement the shell command recipe to query `gcloud compute regions list | tail -n +2 | awk '{print $$1}'` and sort results
    - [ ] Format output into Go string slice and safely update `internal/run/api/region/region.go`
    - [ ] Run `go fmt ./internal/run/api/region/...` post-update
    - [ ] Register `regions-update` in `.PHONY` and add help comment
- [ ] Task: Conductor - User Manual Verification 'Makefile Implementation & Region Update Command' (Protocol in workflow.md)

## Phase 3: Verification & Documentation
- [ ] Task: Verify Codebase Quality and Generated Regions
    - [ ] Run `make regions-update` (if gcloud is available) or verify with test inputs
    - [ ] Run `make test`, `make lint`, and `make build` to guarantee compliance
- [ ] Task: Document Command in README.md
    - [ ] Add `make regions-update` explanation under development commands in `README.md`
- [ ] Task: Conductor - User Manual Verification 'Verification & Documentation' (Protocol in workflow.md)
