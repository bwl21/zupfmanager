# github.com/bwl21/zupfmanager

A desktop application built with Go, Wails, Vue, and ENT ORM with SQLite.

## Project Structure

- **Frontend**: Vue.js with Vite
- **Backend**: Go with Wails
- **Database**: SQLite with ENT ORM

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

## Database

The application uses SQLite with ENT ORM for data persistence. The database schema is defined in the `internal/ent` directory.
