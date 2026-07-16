## 2026-05-22 - [GCP Client Caching]
**Learning:** Initializing Google Cloud Platform clients on every request introduces significant latency (~300ms+) due to credential discovery and gRPC connection overhead. This is particularly problematic in TUI applications that perform concurrent regional requests.
**Action:** Use a thread-safe lazy-initialization pattern with `sync.Mutex` and `context.Background()` for credential discovery to cache and reuse the client across the application lifecycle.
