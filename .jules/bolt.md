## 2026-03-05 - GCP Client Initialization Bottleneck
**Learning:** Initializing a new GCP client (including credential discovery and connection establishment) on every API request introduces significant latency, especially during concurrent multi-region operations. Caching the client and credentials using a thread-safe lazy-initialization pattern dramatically reduces this overhead.
**Action:** Use a `sync.Mutex` and a dedicated `getInternalClient` helper to cache and reuse GCP clients within implementation structs.
