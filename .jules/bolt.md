# Bolt's Journal

## 2026-03-09 - GCP ID Token Credentials Caching & Context Stability
**Learning:** Discovered that every invocation of `GetIDToken` triggered full Google application-default credentials discovery. This involves file I/O (credential files parsing) and GCP Metadata Server HTTP checks, adding up to ~300ms overhead on every ID token request. Caching `google.Credentials` using a thread-safe pattern (`sync.Mutex` with a package-level variable) completely eliminates this overhead on subsequent calls. Crucially, passing `context.Background()` during credentials discovery prevents transient, canceled request contexts from invalidating or closing the cached credentials' network transports.
**Action:** Always lazily initialize and cache GCP default credentials (`google.Credentials`) inside identity token clients using a thread-safe mutex and robust context initialization.

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
