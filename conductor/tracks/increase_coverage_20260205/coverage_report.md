# Coverage Assessment Report

**Date:** 2026-02-05
**Overall Coverage:** 75.0%

## Low Coverage Areas (<80%)

### Critical Infrastructure
- `cmd/run`: 0.0% (Entry point)
- `internal/run/auth`: 69.1%
- `internal/run/command`: 55.6%

### API Clients
- `internal/run/api/service`: 76.0%
  - `UpdateAuthentication`: 0.0%
  - `GetService`: 77.8%

### TUI Applications
- `internal/run/tui/app`: 53.8% (Main app logic)
- `internal/run/tui/app/service/auth`: 66.0%
- `internal/run/tui/app/service/revision`: 55.1%
- `internal/run/tui/app/service/scale`: 64.9%
- `internal/run/tui/app/service/traffic`: 9.2% (Very low)
- `internal/run/tui/app/workerpool/scale`: 67.8%

### TUI Components
- `internal/run/tui/component/footer`: 0.0%
- `internal/run/tui/component/table`: 47.4%

### Models (Missing Tests)
The following packages have no test files:
- `internal/run/model/common/condition`
- `internal/run/model/common/container`
- `internal/run/model/common/env`
- `internal/run/model/common/info`
- `internal/run/model/common/keytopath`
- `internal/run/model/common/project`
- `internal/run/model/common/resources`
- `internal/run/model/common/secret`
- `internal/run/model/common/volume`
- `internal/run/model/domainmapping`
- `internal/run/model/job`
- `internal/run/model/job/execution`
- `internal/run/model/service`
- `internal/run/model/service/networking`
- `internal/run/model/service/scaling`
- `internal/run/model/service/security`
- `internal/run/model/service/traffic`
- `internal/run/model/workerpool`
- `internal/run/model/workerpool/scaling`
