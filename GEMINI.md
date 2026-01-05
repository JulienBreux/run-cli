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
    *   `internal/run/tui/`: Contains the terminal UI logic.
        *   `app/`: Main application layout and pages (`service`, `job`, `worker`).
        *   `component/`: Reusable UI components (`header`, `table`, `loader`, `spinner`, `logo`, `modals`).
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

*   **Project Structure:** Follows standard Go project layout with `cmd` for executables and `internal` for private library code. The `pkg` is used for public libraries.
*   **Linting:** `golangci-lint` is used to enforce code style and quality. Ensure `make lint` passes before committing.
*   **Testing:** Unit tests are encouraged. Use `make test` to verify changes.
*   **UI Components:** The UI is modularized in `internal/run/tui`.
    *   **Pages:** Located in `internal/run/tui/app`.
    *   **Components:** Shared elements (Spinner, Logo, Header) are in `internal/run/tui/component`.

## Application Flow

*   **Initialization:** The application performs a parallel pre-load of **Projects** and **Services** during the initial loading screen (featuring the Logo and Spinner). This ensures the dashboard is populated immediately upon launch.
*   **Async Operations:** Long-running tasks like scaling services or worker pools use the `spinner` component to provide visual feedback without blocking the UI.
*   **Data Fetching:** UI packages (`service`, `project`) expose helper functions (`Fetch`, `PreLoad`) to allow data retrieval without tight coupling to the main application logic.
