# Bolt's Journal

## 2026-03-09 - Lock-Free Pre-Allocated Map-Reduce for Multi-Regional Queries
**Learning:** Performing multi-regional queries (like `listAllRegions` across 24 GCP regions) concurrently using a shared slice protected by `sync.Mutex` causes lock contention, slice allocation churn, and heavy garbage collection overhead. Utilizing a lock-free, pre-allocated map-reduce pattern using a slice of slices (`[][]T`) eliminates mutex contention entirely, and pre-allocating the final result slice based on the exact sum of sizes of individual regional results reduces heap allocations considerably.
**Action:** Always favor lock-free, index-mapped pre-allocation arrays (`[][]T`) over shared mutex-protected slices (`[]T`) when collecting results concurrently from a known, fixed number of concurrent tasks.

## 2026-03-08 - GCP Logging Client Caching & Connection Longevity
**Learning:** Establishing the GCP Stackdriver Logging client requires repeated Google credential discovery and connection establishment, causing high latency (~300ms) inside a reactive TUI interface. Caching `logadmin.Client` instances via a project-aware map with thread-safe `sync.Mutex` ensures subsequent streaming and log extraction operations are instantaneous. Crucially, calling `Close()` on individual stream terminations must be a no-op to prevent premature teardown of connection pools shared across other active streaming views.
**Action:** Keep GCP Logging clients cached globally by project and handle connection termination via a no-op `Close` method, while adding test-isolation resets in unit tests.

## 2026-03-07 - GCP Execution & Project Client Caching
**Learning:** Incomplete caching of GCP client wrappers in remaining packages (`job/execution` and `project`) meant that loading job execution tables or searching for GCP projects still suffered from repetitive credential discovery and connection establishment, degrading TUI interactivity. Standardizing stateful `GCPClient` caching via thread-safe lazy-initialization with `sync.Mutex` completely eliminates this overhead.
**Action:** Ensure all Cloud SDK APIs used in the application leverage the stateful lazy-initialization cached client pattern with proper thread-safety.

## 2026-03-06 - GCP Revisions Client Caching
**Learning:** Failing to cache the GCP Revisions client meant that navigating and reloading service revisions in the TUI invoked credential discovery and gRPC setup on every revision fetch. Standardizing the stateful `sync.Mutex` lazy-initialization pattern ensures that once a client is cached, subsequent service revision lists are lightning fast.
**Action:** Consistently inspect all subcommand and subpackage API wrappers to ensure that client connections are pooled and not destroyed on request boundaries.

## 2026-03-05 - GCP Services Client Caching
**Learning:** Initializing Google Cloud Platform API clients and discovering credentials (ADC) on every single request in a terminal user interface (TUI) introduces a major latency bottleneck (~300ms per request) due to repetitive file system lookups and TCP/gRPC handshakes. Caching the client wrapper using a thread-safe pattern (`sync.Mutex`) avoids this completely. Importantly, client construction must use `context.Background()` rather than request-scoped contexts to prevent cancellation of the shared client connection pool when a single request context is canceled.
**Action:** Always verify if external API or cloud clients are lazily initialized and cached as singletons/long-lived clients across successive TUI interactions, rather than being created and closed on every API call.

## 2025-05-15 - [GCP Client Creation Overhead]
**Learning:** Creating a new GCP client for every API call (especially in a TUI that frequently refreshes and can list all regions) introduces significant latency due to repeated credential discovery, TLS handshakes, and gRPC connection establishment. When listing "All" regions, this results in 24 simultaneous client creations and connection setups.
**Action:** Reuse a single, thread-safe GCP client per package (service, job, etc.) to leverage connection pooling and reduce overhead.
