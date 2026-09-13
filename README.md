# Cloud Run Interactive CLI

<p align="center">
  <a href="https://github.com/JulienBreux/run-cli" target="_blank"><img src="docs/assets/run.gif" /></a>
</p>

**Run CLI** is an interactive CLI to manage your Google Cloud Run resources in a visual way.

The idea is to establish a high-performance, terminal-based alternative to the Google Cloud Console for executing daily operational tasks.

[![Go version](https://img.shields.io/github/go-mod/go-version/JulienBreux/run-cli)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/JulienBreux/run-cli)](https://goreportcard.com/report/github.com/JulienBreux/run-cli)
[![codecov](https://codecov.io/github/JulienBreux/run-cli/graph/badge.svg?token=NSaxKHzR64)](https://codecov.io/github/JulienBreux/run-cli)
[![Release](https://img.shields.io/github/release/JulienBreux/run-cli.svg?style=flat-square)](https://github.com/JulienBreux/run-cli/releases)
[![LICENSE](https://img.shields.io/github/license/JulienBreux/run-cli)](https://github.com/JulienBreux/run-cli/blob/main/LICENSE)

</div>

## ⭐️ Value Proposition

- **Visual Resource Management:** Offers a sophisticated visual representation of cloud resources using interactive tables and dashboards, providing superior clarity over standard flat CLI output.
- **Interactive Discoverability:** Simplifies the management of complex settings through intuitive interactive menus, eliminating the need to memorize extensive and complex command-line flags.
- **Rapid Context Switching:** Enables instantaneous switching between different Google Cloud projects and regions within a single, continuous terminal session.


## ✨ Features

### 🌍 Global

*   **Interactive TUI:** A user-friendly terminal interface to manage your Cloud Run resources.
*   **Project & Region Selection:** Easily switch between your Google Cloud projects and regions.
*   **Log Viewer:** Stream logs from your services directly in the terminal.
*   **Konami Code:** Try the legendary code for a little surprise!

### 🚀 Services

*   **Service Management:** View, search, and manage your Cloud Run services.
*   **Authentication Management:** Toggle between "Require authentication" and "Allow unauthenticated invocations" directly from the interface.
*   **Service Dashboard:** Navigate to a dedicated dashboard for each service with multiple views.
*   **Service Proxy:** Locally proxy private services with automatic authentication token injection (press `p`).
*   **Networking View:** Monitor ingress settings, endpoints status (URI, IAP), and VPC Access configurations.
*   **Security View:** Check authentication requirements, service identity, encryption keys, and binary authorization policies.
*   **Revision Management:** Detailed list of revisions with traffic allocation, tags, and deployment history.
*   **Deep Insights:** Explore revision details including billing mode, startup CPU boost, concurrency, and request timeouts.
*   **Resource Monitoring:** View container configurations, images, ports, and resource limits (Memory, CPU, and GPU/Accelerators).

### ⚡ Jobs

*   **Job Management:** Monitor and manage your Cloud Run jobs.
*   **Job Dashboard:** Dedicated view for jobs including execution history and status.
*   **Execution Management:** View detailed execution history with task success/failure counts, duration, and status.

### 🖥️ Instances

*   **Instance Management:** View, inspect, and manage your Google Cloud Run instances (`Ctrl+N`).
*   **Direct Lifecycle Actions:** Start (`s`), Stop (`x`), and Delete (`k`) instances interactively.
*   **Authentication Management:** Toggle between "Require authentication" and "Allow unauthenticated invocations" directly from the interface (`a`).
*   **Instance Proxy:** Locally proxy private instances with automatic authentication token injection (`p`) and open the active URL in your default browser (`o`).
*   **Instance Dashboard:** Explore multi-tab inspection dashboards including Overview, Containers, Conditions & Networking, and raw syntax-highlighted YAML/JSON definitions (`Enter`).
*   **Log Viewer:** Stream Cloud Logging records specific to instances (`l`).
*   **Console Integration:** Direct shortcut (`Ctrl+Z`) to open the selected instance or instances list directly in Google Cloud Console.

### 👷 Worker Pools

*   **Worker Pool Management:** View and manage your Cloud Run worker pools.
*   **Scaling Control:** Monitor and adjust scaling settings.

### 🌐 Domain Mappings

*   **Domain Management:** View your custom domain mappings.
*   **DNS Configuration:** Quickly access DNS record instructions for easy setup.

## 🚀 Installation

Run CLI is available on Linux, OSX and Windows platforms.

* Binaries for Mac OS, Linux and Windows are available as tarballs in the [release](https://github.com/JulienBreux/run-cli/releases) page.

* Via Homebrew (Mac OS) or LinuxBrew (Linux)

   ```shell
   brew tap julienbreux/run
   brew install --cask julienbreux/run/run
   ```

* Via `go get`

    You can install **Run CLI** using `go install`:

    ```shell
    go install github.com/JulienBreux/run-cli/cmd/run@latest
    ```

* Via CURL

    You can install **Run CLI** using CURL and the shell script:

    ```shell
    curl -sL https://JulienBreux.github.io/run-cli/get | sh
    ```

## ▶️ Usage

Simply run the command:

```sh
gcloud auth application-default login # Currently mandatory :/
run
```

This will start the interactive TUI, allowing you to manage your Google Cloud Run resources.

## 🛠️ Development

This project uses a `Makefile` to streamline development.

### Prerequisites

*   [Go](https://go.dev/doc/install) (version 1.27+)
*   [Docker](https://docs.docker.com/get-docker/) (for building the Docker image)
*   [Google Cloud SDK](https://cloud.google.com/sdk) (`gcloud`, optional, for updating the regions list)

### Build from source

To build the binary from source, run:

```sh
make build
```

This will create the `run` executable in the `./bin` directory.

### Running tests

To run the unit tests:

```sh
make test
```

### Linting

To lint the codebase:

```sh
make lint
```

### Updating Regions

To update the list of supported regions in `internal/run/api/region/region.go` using the `gcloud` CLI:

```sh
make regions-update
```

## 💪 Contributing

Contributions are welcome! We appreciate all contributions, from bug reports and documentation enhancements to code and feature additions.

- **Contribution Guidelines**: Please review our [Contributing Guidelines](CONTRIBUTING.md) to get started with setup, development targets, and conventions.
- **Maintainers & Contributors**: See our roster of maintainers and contributors in [MAINTAINERS.md](MAINTAINERS.md).
- **Code of Conduct**: This project is governed by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).
- **Support & Questions**: Need help or want to discuss ideas? Refer to [SUPPORT.md](SUPPORT.md).
- **Security**: To report a security vulnerability, please see [SECURITY.md](SECURITY.md).

## 📄 License

This project is licensed under the Apache 2.0 License. See the [LICENSE](https://github.com/JulienBreux/run-cli/blob/main/LICENSE) file for details.
