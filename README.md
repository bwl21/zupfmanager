# Zupfmanager

Zupfmanager is a specialized tool for managing and building music projects for zither (zupf) instruments. It helps organize ABC notation music files, manage project configurations, and generate sheet music using the Zupfnoter renderer.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Features
# Features - This section describes the main features of the Zupfmanager.

- Manage music projects with customizable configurations
- Organize songs in projects with priority and difficulty settings
- Import ABC notation files into the song database
- Build projects to generate sheet music in PDF format
- Command-line interface (CLI) for automation and scripting
- Terminal user interface (TUI) for interactive usage

## Installation
# Installation - This section describes how to install the Zupfmanager.

### Prerequisites

- Go 1.23 or higher
- Node.js (for Zupfnoter rendering)

### From Source

```bash
# Clone the repository
git clone https://github.com/bwl21/zupfmanager.git
cd zupfmanager

# Build the application
make build
```

### Binary Releases

Download the latest binary for your platform from the [Releases](https://github.com/bwl21/zupfmanager/releases) page.

## Usage
# Usage - This section describes how to use the Zupfmanager.

### Command Line Interface

Zupfmanager provides a comprehensive CLI for managing projects and songs:

```bash
# Show help
zupfmanager --help

# Project Management
zupfmanager project list                                # List all projects
zupfmanager project create "My Project" "MP"            # Create a new project
zupfmanager project show <project-id>                   # Show project details
zupfmanager project edit-config <project-id>            # Edit project configuration
zupfmanager project add-song <project-id> <song-id>     # Add a song to a project
zupfmanager project remove-song <project-id> <song-id>  # Remove a song from a project
zupfmanager project edit-song <project-id> <song-id>    # Edit song settings in a project
zupfmanager project build <project-id>                  # Build a project

# Song Management
zupfmanager song list                                   # List all songs
zupfmanager song show <song-id>                         # Show song details
zupfmanager song search <query>                         # Search for songs

# Import ABC Files
zupfmanager import <directory>                          # Import ABC files from a directory

# Interactive Terminal UI
zupfmanager ui                                          # Launch the terminal UI
```

### Terminal User Interface

The application includes an interactive Terminal UI that provides the same functionality in a more user-friendly way:

```bash
# Launch the TUI
zupfmanager ui
```

Navigate the interface using keyboard shortcuts:
- Arrow keys or j/k: Navigate lists
- Enter: Select/View item
- q: Quit
- Esc: Go back
- a: Add item
- e: Edit item
- d: Delete item

### Building Projects

To build a project and generate sheet music:

```bash
zupfmanager project build <project-id> --abc-file-dir <abc-files-directory> --output-dir <output-directory>
```

This will:
1. Process all songs in the project
2. Merge project configuration with song-specific configurations
3. Generate PDF sheet music using Zupfnoter
4. Save output files to the specified directory

## Data Model
# Data Model - This section describes the data model used by Zupfmanager.

Zupfmanager uses a simple data model with three main entities:

1. **Project**: A collection of songs with a title, short name, and configuration
2. **Song**: An ABC notation file with metadata
3. **ProjectSong**: A join entity that connects projects and songs with additional attributes like priority and difficulty

## Configuration
# Configuration - This section describes how to configure projects.

Projects can be configured with custom settings that override or extend the default configurations in ABC files. The configuration is stored in JSON format and can be edited using the `project edit-config` command.

## ABC File Format
# ABC File Format - This section describes the ABC file format used by Zupfmanager.

Zupfmanager works with ABC notation files that include Zupfnoter-specific configuration. A typical ABC file structure:

```
X:123
T:Song Title
...ABC notation...

%%%%zupfnoter.config
{
  "produce": [1, 2],
  "extract": {
    "0": {
      "title": "View",
      ...
    },
    "1": {
      "title": "Extract 1",
      ...
    }
  }
}
```

## Development
# Development - This section describes how to set up a development environment, build, run, and debug the application.

### Prerequisites

- Go 1.23 or higher
- Node.js (for Zupfnoter rendering)
- SQLite (for local development database)
- Make (for build automation)

### Project Structure

- **cmd/**: Command-line interface implementation using Cobra
- **internal/**: Internal packages
  - **database/**: Database connection and management
  - **ent/**: Entity definitions and database schema
  - **ui/**: Terminal user interface using Bubble Tea
  - **zupfnoter/**: Integration with Zupfnoter renderer
- **testdata/**: Sample ABC files for testing

### Setting Up Development Environment

1. **Clone the repository**:
   ```bash
   git clone https://github.com/bwl21/zupfmanager.git
   cd zupfmanager
   ```

2. **Install dependencies**:
   ```bash
   # Install Go dependencies
   go mod download
   
   # Install Zupfnoter dependencies
   cd internal/zupfnoter
   npm install
   cd ../..
   ```

### Building the Application

```bash
# Build for current platform (outputs to bin/zupfmanager)
make build

# Build for all supported platforms
make release

# Install to $GOPATH/bin
make install
```

### Running the Application

```bash
# Run the CLI
./bin/zupfmanager --help

# Run the TUI
./bin/zupfmanager ui

# Run with debug logging
DEBUG=1 ./bin/zupfmanager --log-level=debug
```

### Debugging

#### Using Delve (recommended)

1. **Install Delve**:
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   ```

2. **Debug with Delve**:
   ```bash
   # Start debugger
   dlv debug . -- [command] [flags]
   
   # Example: Debug project build
   dlv debug . -- project build 1 --abc-file-dir ./testdata --output-dir ./output
   ```

   Common Delve commands:
   - `break main.main`: Set breakpoint at main function
   - `break pkg.FunctionName`: Set breakpoint at specific function
   - `continue`: Continue execution
   - `next`: Step over to next line
   - `step`: Step into function
   - `print variable`: Print variable value
   - `exit`: Quit debugger

#### Using GoLand

1. **Setting Up Debug Configuration**
   - Click on "Add Configuration..." in the top toolbar
   - Click "+" and select "Go Build"
   - Configure the run configuration:
     - Name: `Debug Zupfmanager`
     - Run kind: `Directory`
     - Directory: Select the project root directory
     - Working directory: Select the project root directory
     - Program arguments: `project build 1 --abc-file-dir ./testdata --output-dir ./output`
     - Environment: `DEBUG=1`

2. **Using the Debugger**
   - Set breakpoints by clicking in the gutter next to the line numbers
   - Start debugging by clicking the bug icon or pressing `Shift + F9`
   - Use the debug toolbar to:
     - Step over (F8)
     - Step into (F7)
     - Step out (Shift + F8)
     - Resume program (F9)
     - Stop (Ctrl+F2)

3. **Debugging Features**
   - **Watches**: Add variables to watch in the "Watches" panel
   - **Variables**: Inspect all variables in the current scope
   - **Debug Console**: Execute code in the current context
   - **Evaluate Expression**: Select code and press Alt+F8 to evaluate
   - **Conditional Breakpoints**: Right-click a breakpoint to add conditions

4. **Remote Debugging**
   - Create a new "Go Remote" configuration
   - Start the application with Delve:
     ```bash
     dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient . -- project build 1
     ```
   - Connect from GoLand using the remote configuration

5. **Debugging Tests**
   - Right-click on a test file or test function
   - Select "Debug 'Test...'"
   - Use the same debugging features as with the main application

#### Using VS Code

1. Install the Go extension for VS Code
2. Create or update `.vscode/launch.json`:
   ```json
   {
       "version": "0.2.0",
       "configurations": [
           {
               "name": "Debug Zupfmanager",
               "type": "go",
               "request": "launch",
               "mode": "auto",
               "program": "${workspaceFolder}",
               "args": ["project", "build", "1", "--abc-file-dir", "./testdata", "--output-dir", "./output"],
               "env": {
                   "DEBUG": "1"
               },
               "showLog": true
           }
       ]
   }
   ```
3. Set breakpoints in your code and press F5 to start debugging

### Testing

```bash
# Run all tests
make test

# Run specific test
go test -v ./... -run TestFunctionName

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
```

### Database Management

```bash
# Generate database models (after schema changes)
make generate

# View database (requires sqlite3)
sqlite3 data/zupfmanager.db
```

### Common Development Tasks

#### Adding a New Command

1. Create a new file in `cmd/` (e.g., `cmd/mycommand.go`)
2. Define your command using Cobra
3. Register it in `cmd/root.go`
4. Run `make generate` to update documentation

#### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/example/package

# Tidy up go.mod
go mod tidy
```

### Troubleshooting

- **Permission denied errors**: Try running with `sudo` or fix file permissions
- **Database issues**: Delete `data/zupfmanager.db` and restart the application
- **Zupfnoter not found**: Make sure Node.js is installed and Zupfnoter dependencies are installed
- **Build failures**: Run `go clean -modcache` and try again

### Code Style

- Follow standard Go formatting: `gofmt -s -w .`
- Run `make lint` to check for style issues
- Document all exported functions and types
- Write tests for new functionality

## License
# License - This section describes the license for the Zupfmanager.

This project is licensed under the MIT License - see the LICENSE file for details.
