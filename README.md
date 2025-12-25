<div align="center">

```
 ____  _   _ _   _
|  _ \| | | | \ | |
| |_) | | | |  \| |
|  _ <| |_| | . ` |
|_| \_\_____|_| \_|
```

# Run CLI

**Run CLI** is an interactive CLI to manage your Google Cloud Run resources with panache!

[![Go version](https://img.shields.io/github/go-mod/go-version/JulienBreux/run-cli)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/JulienBreux/run-cli)](https://goreportcard.com/report/github.com/JulienBreux/run-cli)
[![LICENSE](https://img.shields.io/github/license/JulienBreux/run-cli)](https://github.com/JulienBreux/run-cli/blob/main/LICENSE)

</div>

## ✨ Features

*   **Interactive TUI:** A user-friendly terminal interface to manage your Cloud Run resources.
*   **Service Management:** View, search, and manage your Cloud Run services.
*   **Job Management:** Monitor and manage your Cloud Run jobs.
*   **Log Viewer:** Stream logs from your services directly in the terminal.
*   **Project & Region Selection:** Easily switch between your Google Cloud projects and regions.

## 🚀 Installation

Run CLI is available on Linux, OSX and Windows platforms.

* Binaries for Mac OS, Linux and Windows are available as tarballs in the [release](https://github.com/JulienBreux/run-cli/releases) page.

* Via Homebrew (Mac OS) or LinuxBrew (Linux)

   ```shell
   brew tap julienbreux/tap/run
   brew install --cask julienbreux/tap/run
   ```

* Via `go get`

    You can install **Run CLI** using `go install`:

    ```shell
    go install github.com/JulienBreux/run-cli/cmd/run@latest
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

*   [Go](https://go.dev/doc/install)
*   [Docker](https://docs.docker.com/get-docker/) (for building the Docker image)

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

## 💪 Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## 📄 License

This project is licensed under the Apache 2.0 License. See the [LICENSE](https://github.com/JulienBreux/run-cli/blob/main/LICENSE) file for details.
