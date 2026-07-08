# Bolt's Journal

## 2025-05-15 - [GCP Client Creation Overhead]
**Learning:** Creating a new GCP client for every API call (especially in a TUI that frequently refreshes and can list all regions) introduces significant latency due to repeated credential discovery, TLS handshakes, and gRPC connection establishment. When listing "All" regions, this results in 24 simultaneous client creations and connection setups.
**Action:** Reuse a single, thread-safe GCP client per package (service, job, etc.) to leverage connection pooling and reduce overhead.
