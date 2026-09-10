# Specification: Add Makefile Command to Update Cloud Run Regions

## Overview
Add a new `make regions-update` target in the project `Makefile` that queries Google Cloud for available compute regions using the `gcloud` CLI (`gcloud compute regions list | tail -n +2 | awk '{print $1}'`), sorts and formats them into Go string slice syntax, and updates the `List()` function in `internal/run/api/region/region.go`.

## Functional Requirements
1. **Makefile Target (`regions-update`)**:
   - Add a new target `regions-update: ## Update regions in region.go using gcloud` to `Makefile`.
   - Add `regions-update` to `.PHONY`.
   - Validate that `gcloud` is installed and available in `PATH`; fail with an informative error if absent.
   - Execute `gcloud compute regions list | tail -n +2 | awk '{print $$1}' | sort` to fetch regions.
   - Safely update `internal/run/api/region/region.go` with the fetched regions, preserving the license header, package declaration, and comments.
   - Run `go fmt ./internal/run/api/region/...` to ensure clean formatting.
2. **Documentation & Tests**:
   - Update `README.md` to document the new `make regions-update` command under development/build instructions.
   - Verify that unit tests in `internal/run/api/region/region_test.go` pass and cover the region list functionality.

## Non-Functional Requirements
- **Determinism**: Region output must be sorted alphabetically and formatted cleanly using standard Go code style (`go fmt`).
- **Error Handling**: Fail fast with clear user feedback if `gcloud` is missing or fails.
- **Safety**: Do not overwrite or corrupt `region.go` if the `gcloud` command fails or returns empty output.

## Acceptance Criteria
- [ ] Running `make help` displays the `regions-update` command with its description.
- [ ] Running `make regions-update` checks for `gcloud`, fetches the regions list, formats them, updates `internal/run/api/region/region.go`, and formats the file.
- [ ] `internal/run/api/region/region.go` maintains its license header, `ALL` constant, and updated `List()` function.
- [ ] `make test`, `make lint`, and `make build` pass cleanly.
- [ ] `README.md` documents `make regions-update`.

## Out of Scope
- Automatic periodic execution via CI (the command is intended as a manual maintenance command for developers).
- Fetching regions directly at runtime via Google Cloud Run API client (hardcoded regions list is intentional for offline fallback / rapid local UI response).
