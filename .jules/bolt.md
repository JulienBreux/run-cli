## 2026-05-27 - GCP Client Initialization Overhead
**Learning:** Creating a new GCP gRPC client and discovering credentials for every single request adds significant latency (hundreds of milliseconds) due to credential discovery (especially in metadata server environments) and TLS/gRPC handshake overhead.
**Action:** Always implement lazy-initialization and caching for GCP clients in long-lived components or TUI applications where multiple sequential or concurrent calls are expected. Use a sync.Mutex for thread-safe access to the shared client.

## 2026-07-14 - GoReleaser CI in Pull Requests
**Learning:** Running GoReleaser in a pull request environment without the --snapshot flag causes failures because the current commit doesn't match the latest git tag.
**Action:** Use conditional logic in GitHub Actions to add the --snapshot flag to goreleaser-action when triggered by a pull_request event.
