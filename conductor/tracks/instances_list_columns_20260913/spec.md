# Specification: Instance List Columns and Auth Styling

## Overview
Update the Cloud Run Instances list table in the TUI to align its visual structure with Services. The `AUTH` column remains as the first titled column (at column index 1, immediately following the proxy indicator at index 0) displaying colored Yes/No status (`[red]Yes` for required authentication, `[green]No` for unauthenticated allowed). The `CONTAINERS` column is removed from the table view, and column expansion weights are adjusted to provide additional horizontal space for instance names and timestamps.

## Functional Requirements
1. **Remove CONTAINERS Column**:
   - Remove `"CONTAINERS"` header and its corresponding table cell from `internal/run/tui/app/instance/instance.go`.
   - Free up width allocation formerly used by containers.
2. **Column Arrangement and Header Alignment**:
   - Align columns to:
     - Index 0: `""` (Proxy status indicator `[green]P`)
     - Index 1: `"AUTH"` (Authentication status)
     - Index 2: `"NAME"` (Instance ID/name)
     - Index 3: `"REGION"` (Deployment region)
     - Index 4: `"STATUS"` (Instance status)
     - Index 5: `"LAST UPDATED"` (Relative update timestamp)
3. **Authentication Value Styling**:
   - Render `[red]Yes` when `InvokerIamDisabled` is false (authentication required).
   - Render `[green]No` when `InvokerIamDisabled` is true (unauthenticated invocations permitted), consistent with the Services table.
4. **Column Expansions**:
   - Update expansion weights to `[]int{1, 1, 3, 1, 1, 2}` to optimize horizontal space distribution.

## Non-Functional Requirements
- **Consistency**: Maintain visual parity with the Services list table layout.
- **Test Coverage**: Update all unit tests in `internal/run/tui/app/instance/instance_test.go` to validate new column indexes, values, and total column counts. Maintain 100% statement coverage.
- **Linting**: Ensure code passes `make lint` without errors.

## Acceptance Criteria
- [ ] Instances table displays 6 columns: `""`, `"AUTH"`, `"NAME"`, `"REGION"`, `"STATUS"`, `"LAST UPDATED"`.
- [ ] No `CONTAINERS` column is displayed in the Instances list.
- [ ] Auth values are displayed with `[red]Yes` and `[green]No`.
- [ ] All unit tests pass (`make test`).
- [ ] Linter passes (`make lint`).
- [ ] Binary builds successfully (`make build`).

## Out of Scope
- Modifications to instance detail dashboard (`dashboard.go`) where container information continues to be available.
- Alterations to API fetching or GCP backend queries.
