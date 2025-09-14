# Bow

Bow is an interactive terminal user interface (TUI) utility for managing Arcanist diffs in Phabricator. It simplifies the process of selecting diffs, commits, and base commits for creating or updating diffs.

## Features

- **Interactive Diff Selection**: Browse and select existing diffs from Arcanist.
- **Commit Selection**: Choose commits to diff on and base commits to diff from.
- **Command Execution**: Create new diffs or update existing ones.
- **Terminal UI**: User-friendly text-based interface for navigation and selection.
- **Git Integration**: Fetches commit history from the local Git repository.

## Requirements

- Go 1.25.1 or later
- Arcanist (Phabricator's command-line tool)
- A Git repository with commits

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bow.git
   cd bow
   ```

2. Build the application:
   ```bash
   go build
   ```

   Or install globally:
   ```bash
   go install
   ```

## Usage

Run Bow in a Git repository directory:

```bash
./bow
```

The TUI will display panels for:
- **Diff from**: Select the base commit
- **Diff on**: Select the target commit
- **Diff to update**: Choose an existing diff
- **Command**: Select "Create" or "Update" action

Use arrow keys to navigate, Enter to select, and follow on-screen instructions.
