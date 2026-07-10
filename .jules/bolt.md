## 2026-03-01 - Caching GCP Clients and Credentials
**Learning:** Initializing a new GCP client (especially gRPC-based ones like Cloud Run v2) for every API call is extremely expensive due to credential discovery and connection handshakes. This is particularly noticeable in TUIs making concurrent requests for multiple regions.
**Action:** Use a `sync.Mutex` and a private field in the API client struct to lazily initialize and cache the client, ensuring thread-safe reuse of connections and credentials across the application's lifecycle.
