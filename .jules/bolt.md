## 2026-05-26 - Optimized GCP client reuse across all API packages
**Learning:** Initializing Google Cloud Platform (GCP) API clients and fetching credentials for every API call introduces significant latency due to disk I/O, network handshakes, and gRPC connection setup. Caching these clients using a lazy-initialization pattern with `sync.Mutex` significantly improves performance in a concurrent environment like a TUI.
**Action:** Use thread-safe client caching for all long-lived API clients. Ensure unit tests reset any global state to maintain isolation.
