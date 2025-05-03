# github.com/bwl21/zupfmanager

A desktop application built with Go, Wails, Vue, and ENT ORM with SQLite.

## Project Structure

- **Frontend**: Vue.js with Vite
- **Backend**: Go with Wails
- **Database**: SQLite with ENT ORM
- **CLI**: Command Line Interface using Cobra
- **TUI**: Terminal User Interface using Bubble Tea

## Development

To run in live development mode:

```bash
wails dev
```

This will run a Vite development server for the frontend with hot reload. You can also access the Go methods through the dev server at http://localhost:34115.

## Building

To build a production package:

```bash
wails build
```

## Command Line Interface

The application provides a CLI for managing projects and songs:

```bash
# Show help
zupfmanager --help

# List projects
zupfmanager project list

# List songs
zupfmanager song list

# Show project details
zupfmanager project show <project-id>

# Add a song to a project
zupfmanager project add <project-id> <song-id> --priority 1 --difficulty medium
```

## Terminal User Interface

The application also includes an interactive Terminal UI that provides the same functionality in a more interactive way:

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

## Database

The application uses SQLite with ENT ORM for data persistence. The database schema is defined in the `internal/ent` directory.
