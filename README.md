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
# Development - This section describes the project structure and how to build from source.

### Project Structure

- **cmd/**: Command-line interface implementation using Cobra
- **internal/**: Internal packages
  - **database/**: Database connection and management
  - **ent/**: Entity definitions and database schema
  - **ui/**: Terminal user interface using Bubble Tea
  - **zupfnoter/**: Integration with Zupfnoter renderer
- **testdata/**: Sample ABC files for testing

### Building from Source

```bash
# Build for current platform
make build

# Build for all supported platforms
make release
```

## License
# License - This section describes the license for the Zupfmanager.

This project is licensed under the MIT License - see the LICENSE file for details.
