# Org Agenda CLI Specification

## Overview
A CLI tool to parse Emacs Org-mode files and manage agendas and TODO lists directly from the terminal.

## Architecture
- **Language**: Go
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Configuration**: YAML format (via [Viper](https://github.com/spf13/viper))

## Command Structure

### Global Options
- `--help (-h)`: Display help information
- `--version (-v)`: Display version information

### Commands

#### 1. `agenda`
Displays the agenda view. Aggregates tasks with schedules and deadlines within a specified period.

- **Usage**: `org-agenda agenda [flags]`
- **Flags**:
    - `--range <day|week>`: Specify the display range (default: `day`).
    - `--date <YYYY-MM-DD>`: Specify the reference date (default: today).
    - `--tag <tag>`: Filter items by a specific tag.
    - `--tui`: Enable interactive TUI mode.

#### 2. `todo`
Manages TODO items.

- **Usage**: `org-agenda todo [command] [flags]`
- **Subcommands**:
    - `list`: Display a list of TODO items (default behavior).
        - `--status <TODO|WAITING|DONE>`: Filter by status.
        - `--tag <tag>`: Filter by tag.
        - `--tui`: Enable interactive TUI mode.
    - `add`: Add a new TODO item.
        - `--file <path>`: Specify the target file (defaults to the configured inbox file).
        - `--schedule <date>`: Set a SCHEDULED timestamp.
        - `--deadline <date>`: Set a DEADLINE timestamp.
        - `--tags <tag1,tag2>`: Set tags.
        - `<title>`: Content of the task.
    - `done`: Mark a task as DONE.
        - `<id|index>`: Specify the task ID or line index.

#### 3. `config`
Manages the configuration file.

- **Usage**: `org-agenda config [command]`
- **Subcommands**:
    - `list`: Display current configuration (e.g., loaded Org file paths).
    - `add-path <path>`: Add an Org file path to the search/display list.

## TUI Interaction
Common keybindings for TUI mode:

### List View
- `j` / `Down`: Move selection down
- `k` / `Up`: Move selection up
- `Enter`: View item details (RawContent)
- `q` / `Esc` / `Ctrl+C`: Quit

### Detail View
- `j` / `Down`: Scroll down
- `k` / `Up`: Scroll up
- `q` / `Esc` / `Backspace`: Return to list view
- `Ctrl+C`: Quit


## Configuration
Configuration follows the XDG Base Directory Specification, typically stored at `~/.config/org-agenda-cli/config.yaml`.

### Config Schema Example
```yaml
org_files:
  - "/home/user/org/work.org"
  - "/home/user/org/private.org"
default_file: "/home/user/org/inbox.org"
```

## Data Model (Go Structs)

### Item
Represents an entry in an Org file.

```go
type Item struct {
    Title       string
    Status      string    // "TODO", "DONE", "WAITING", etc.
    Tags        []string
    Scheduled   *time.Time
    Deadline    *time.Time
    FilePath    string
    LineNumber  int
    RawContent  string    // Body content
}
```

## Code Quality Standards
- **Formatting**: Always use `gofmt` for formatting.
- **Development Process**: Follow Test-Driven Development (TDD).
- **Comments**: Write comments focusing on "why" things are done rather than "what" is done.
