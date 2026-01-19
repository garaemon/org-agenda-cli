# org-agenda-cli

A CLI tool to parse Emacs Org-mode files and manage agendas and TODO lists directly from the terminal.

## Features

- **Agenda View**: Aggregates tasks with schedules and deadlines.
- **TODO Management**: List and add TODO items across multiple Org files.
- **Configurable**: Easily manage the list of Org files to track.
- **Fast**: Built with Go for high performance.

## Installation

### Using go install

```bash
go install github.com/garaemon/devgo@latest
```

### Building from source

You can also build the project from source using `make`:

```bash
make build
# The binary 'org-agenda-cli' will be created in the current directory.
# You can move it to your $PATH, e.g.:
sudo mv org-agenda-cli /usr/local/bin/org-agenda
```

## Usage

### Configuration

Add your Org files to the configuration:

```bash
org-agenda config add-path ~/org/work.org
org-agenda config add-path ~/org/private.org
```

Remove an Org file from the configuration:

```bash
org-agenda config remove-path ~/org/work.org
```

You can manually edit the configuration file (usually `~/.config/org-agenda-cli/config.yaml`) to set up capture templates.

#### Capture Configuration

You can configure where and how notes are captured using the `capture` section in your `config.yaml`.

```yaml
capture:
  # The file to capture to (overrides default_file if set)
  default_file: "/Users/user/org/inbox.org"

  # Format of the captured entry
  # %t: Timestamp <YYYY-MM-DD Mon HH:MM>
  # %c: The content you passed
  # %L: Link to the current working directory
  # %Y: Year (2006)
  # %y: Year (06)
  # %m: Month (01)
  # %d: Day (02)
  # %H: Hour (15)
  # %M: Minute (04)
  # %S: Second (05)
  # %A: Day of week (Monday)
  # %a: Day of week (Mon)
  format: "* %t\n  %c\n  Link: %L"

  # Optional: Insert under a specific heading
  # heading: "Inbox"

  # Optional: Insert under a specific outline path (file+olp)
  # This takes precedence over 'heading'
  # olp:
  #   - "Projects"
  #   - "Random"

  # Optional: Prepend to the file or target heading instead of appending
  prepend: false
```

### Agenda

Display today's agenda:

```bash
org-agenda agenda
```

Display agenda for a week starting from a specific date:

```bash
org-agenda agenda --date 2026-01-01 --range week
```

### TODO List

List all TODO items:

```bash
org-agenda todo list
```

Filter by status or tag:

```bash
org-agenda todo list --status WAITING
org-agenda todo list --tag work
```

### Adding Tasks

Add a new TODO item to the default file:

```bash
org-agenda todo add "Review project proposal" --tags "work,urgent" --schedule 2026-01-05
```

### Capturing Notes

Capture a quick note to your configured Org file:

```bash
org-agenda capture "Check this new tool out"
```

You can override the target file temporarily:

```bash
org-agenda capture "Meeting notes" --file ~/org/meetings.org
```

The capture command respects the `heading`, `olp`, and `prepend` settings in your `config.yaml`. Ideally, this allows you to set up a workflow similar to Emacs `org-capture`.

### Tags

List all unique tags across all configured Org files:

```bash
org-agenda tags
```

## Development

### Prerequisites

- Go 1.23 or later

### Running Tests

```bash
go test ./...
```

## License

MIT
