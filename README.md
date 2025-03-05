# TikTok Live Logger

A beautiful command-line tool to log and view TikTok live streams. Built with Go and featuring a modern TUI interface.

## Features

- Connect to any TikTok live stream
- Log all events (chat messages, gifts, likes, follows, etc.)
- Beautiful TUI interface using Bubble Tea
- Persistent storage using SQLite
- View historical logs with a beautiful interface
- Real-time statistics (viewer count, likes, etc.)
- Configurable settings
- Automatic log cleanup
- Debug mode for troubleshooting

## Installation

```bash
go install github.com/krishna8421/tiktok-live-logger
```

## Usage

### Log a Live Stream

```bash
tiktok-live-logger log username
```

This will:

1. Connect to the specified user's live stream
2. Show a beautiful TUI interface with:
   - Real-time statistics
   - Live chat messages
   - Gifts and other events
3. Save all events to a local SQLite database

### View Saved Logs

```bash
tiktok-live-logger list
```

This will:

1. Show a list of all logged streams
2. Display the number of sessions per stream
3. Allow you to select a stream to view its detailed logs

### Clean Old Logs

```bash
tiktok-live-logger clean [days]
```

This will:

1. Remove logs older than the specified number of days
2. If no days are specified, defaults to 30 days
3. Shows how many events were deleted

### Manage Configuration

```bash
tiktok-live-logger config
```

View current configuration:

```bash
tiktok-live-logger config
```

Update configuration:

```bash
tiktok-live-logger config set <key> <value>
```

Available configuration keys:

- `default_days_to_keep`: Number of days to keep logs (default: 30)
- `database_path`: Path to the SQLite database file
- `debug_mode`: Enable/disable debug mode (true/false)

## Global Options

- `--db, -d`: Specify a custom database path
- `--debug, -v`: Enable debug mode for troubleshooting

## Data Storage

All data is stored in a SQLite database located at:

```bash
~/.tiktok-live-logger/events.db
```

Configuration is stored at:

```bash
~/.tiktok-live-logger/config.json
```

## Keyboard Shortcuts

- `Ctrl+C` or `Esc`: Exit the application
- `↑`/`↓`: Navigate through lists
- `Enter`: Select an item
- `Space`: Scroll through logs

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
