# Bolt's Journal

## 2026-03-08 - GCP Logging Client Caching with Project-Aware Map
**Learning:** Google Cloud Logging client requires a specific `projectID` parameter during initialization. To cache Logging clients effectively in a multi-project TUI environment, standardizing on a project-aware thread-safe map cache with a `sync.Mutex` prevents expensive credential discovery (~300ms) and connection overhead on every log stream, while ensuring different projects resolve to their correct cached clients. Using `context.Background()` avoids connection teardowns due to transient request-level context cancellations, and a no-op `Close()` preserves the shared connections.
**Action:** Use project-aware maps for caching GCP clients that are bound to a specific project ID upon creation, and ensure proper test isolation by resetting global maps.

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
