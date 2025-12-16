# Run CLI Context

## Project Overview

**Run CLI** is an interactive Command Line Interface (CLI) designed to manage Google Cloud Run resources with an enhanced user experience. It leverages a Text-based User Interface (TUI) to provide a rich, interactive environment directly in the terminal.

*   **Language:** Go (Golang)
*   **Key Libraries:**
    *   `github.com/rivo/tview`: For building the Terminal User Interface (TUI).
    *   `github.com/gdamore/tcell`: Underlying library for terminal handling.
    *   `github.com/spf13/cobra`: For CLI command structure and flag parsing.
*   **Architecture:**
    *   `cmd/run/`: Contains the main entry point (`main.go`).
    *   `internal/run/command/`: Handles the CLI command definitions and initialization.
    *   `internal/run/ui/`: Contains the UI logic, including application layout (`app`), pages (`service`, `job`, `worker`), and components (`header`, `table`, `modals`).
    *   `internal/run/model/`: Defines the data models used across the application (e.g., `service`, `job`, `info`).

## Building and Running

The project uses a `Makefile` to automate common development tasks.

### Key Commands

*   **Build:** Compile the project and generate the binary in `./bin/run`.
    ```bash
    make build
    ```

*   **Run:** Build and execute the local binary.
    ```bash
    make run
    ```

*   **Test:** Run unit tests with coverage analysis.
    ```bash
    make test
    ```

*   **Lint:** Run the code linter (`golangci-lint`).
    ```bash
    make lint
    ```

*   **Docker:**
    *   Build the Docker image: `make build-image`
    *   Run the container: `make run-container`

## Development Conventions

*   **Project Structure:** Follows standard Go project layout with `cmd` for executables and `internal` for private library code.
*   **Linting:** `golangci-lint` is used to enforce code style and quality. Ensure `make lint` passes before committing.
*   **Testing:** Unit tests are encouraged. Use `make test` to verify changes.
*   **UI Components:** The UI is modularized in `internal/run/ui`. New pages or components should be added there and registered in `internal/run/ui/app/app.go`.
