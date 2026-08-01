# Bolt's Journal

## 2026-03-09 - Lock-Free Map-Reduce Slice Pattern for Regional Listings
**Learning:** Performing multi-regional operations (such as listing services or domain mappings across all Cloud Run locations) concurrently can lead to mutex lock contention and excessive heap allocation/copy overhead when multiple goroutines write to a shared slice protected by `sync.Mutex`. Utilizing a pre-allocated slice-of-slices pattern indexed by region eliminates the lock entirely. Furthermore, calculating the precise required final slice capacity from the sub-slices and pre-allocating the aggregated slice before flattening avoids intermediate slice growth/copying altogether, improving concurrency and reducing memory allocations.
**Action:** When querying multiple data sources concurrently, use a lock-free slice-of-slices pattern and pre-calculate exact capacity before merging into a single pre-allocated flat slice.

## 2026-03-08 - GCP Logging Client Caching & Connection Longevity
**Learning:** Establishing the GCP Stackdriver Logging client requires repeated Google credential discovery and connection establishment, causing high latency (~300ms) inside a reactive TUI interface. Caching `logadmin.Client` instances via a project-aware map with thread-safe `sync.Mutex` ensures subsequent streaming and log extraction operations are instantaneous. Crucially, calling `Close()` on individual stream terminations must be a no-op to prevent premature teardown of connection pools shared across other active streaming views.
**Action:** Keep GCP Logging clients cached globally by project and handle connection termination via a no-op `Close` method, while adding test-isolation resets in unit tests.

## 2026-03-07 - GCP Execution & Project Client Caching
**Learning:** Incomplete caching of GCP client wrappers in remaining packages (`job/execution` and `project`) meant that loading job execution tables or searching for GCP projects still suffered from repetitive credential discovery and connection establishment, degrading TUI interactivity. Standardizing stateful `GCPClient` caching via thread-safe lazy-initialization with `sync.Mutex` completely eliminates this overhead.
**Action:** Ensure all Cloud SDK APIs used in the application leverage the stateful lazy-initialization cached client pattern with proper thread-safety.

## 2026-03-06 - GCP Revisions Client Caching
**Learning:** Failing to cache the GCP Revisions client meant that navigating and reloading service revisions in the TUI invoked credential discovery and gRPC setup on every revision fetch. Standardizing the stateful `sync.Mutex` lazy-initialization pattern ensures that once a client is cached, subsequent service revision lists are lightning fast.
**Action:** Consistently inspect all subcommand and subpackage API wrappers to ensure that client connections are pooled and not destroyed on request boundaries.

## 2025-05-15 - [GCP Client Creation Overhead]
**Learning:** Creating a new GCP client for every API call (especially in a TUI that frequently refreshes and can list all regions) introduces significant latency due to repeated credential discovery, TLS handshakes, and gRPC connection establishment. When listing "All" regions, this results in 24 simultaneous client creations and connection setups.
**Action:** Reuse a single, thread-safe GCP client per package (service, job, etc.) to leverage connection pooling and reduce overhead.
