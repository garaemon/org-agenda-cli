# org-agenda-cli

A CLI tool to parse Emacs Org-mode files and manage agendas and TODO lists directly from the terminal.

## Features

- **Agenda View**: Aggregates tasks with schedules and deadlines.
- **TODO Management**: List and add TODO items across multiple Org files.
- **Configurable**: Easily manage the list of Org files to track.
- **Fast**: Built with Go for high performance.

## Installation

```bash
go build -o org-agenda main.go
mv org-agenda /usr/local/bin/
```

## Usage

### Configuration

Add your Org files to the configuration:

```bash
org-agenda config add-path ~/org/work.org
org-agenda config add-path ~/org/private.org
```

Set a default file for adding new tasks:

```bash
# Edit ~/.config/org-agenda-cli/config.yaml
default_file: "/home/user/org/inbox.org"
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

## Development

### Prerequisites

- Go 1.23 or later

### Running Tests

```bash
go test ./...
```

## License

MIT
