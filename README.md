# Tempi - Temporary Folder Manager

Tempi is a CLI tool for Windows that creates and manages temporary folders with automatic expiration and cleanup.

## Features

- Create temporary folders with configurable expiration times
- Automatic cleanup of expired folders
- Persistent metadata that survives system reboots
- Track folder activity based on file modification times
- Force delete all temporary folders on demand

## Installation

### Building from Source

```bash
go build -o tempi.exe
```

The binary can be placed anywhere in your PATH for easy access.

## Usage

### Create a New Temporary Folder

Create a folder with default settings (4-hour expiration):

```bash
tempi
```

or explicitly use the `new` command:

```bash
tempi new
```

Create a folder with custom expiration time:

```bash
tempi new --deadtime 2h
tempi new --deadtime 30m
tempi new --deadtime 1h30m
```

### List Active Folders

View all tracked temporary folders and their remaining time:

```bash
tempi list
```

Example output:
```
Active temporary folders (2):

[1] C:\Temp\tempi_20251106_123456 (expires in 3h12m)
[2] C:\Temp\tempi_20251106_133000 (expires in 58m)
```

### Clean Expired Folders

Remove folders that have exceeded their expiration time or are no longer active:

```bash
tempi clean
```

**Note:** Cleanup also runs automatically whenever you execute any Tempi command (except `list` and `help`).

### Delete All Folders Immediately

Force delete all tracked temporary folders regardless of their expiration time:

```bash
tempi deletenow
```

### Help

View help information:

```bash
tempi --help
tempi [command] --help
```

## Configuration

### Default Settings

- **Default Temp Directory**: `C:\Temp`
- **Default Expiration Time**: 4 hours
- **Registry Location**: `%LOCALAPPDATA%\Tempi\registry.json`

### How It Works

1. When you create a temporary folder, Tempi stores its metadata in a registry file at `%LOCALAPPDATA%\Tempi\registry.json`
2. This registry persists across system reboots
3. Folders are automatically cleaned when:
   - Their expiration time has passed
   - No files have been modified within the deadtime period
4. Cleanup happens automatically when you run any Tempi command (opportunistic cleanup)

### Expiration Logic

A folder is considered expired and eligible for cleanup if:
- The expiration time (`created_at + deadtime`) has passed, OR
- The last modification time of any file in the folder exceeds the deadtime

This dual approach ensures that active folders (with recent file modifications) are not prematurely deleted.

## Examples

```bash
# Create a folder for a quick test (expires in 30 minutes)
tempi new --deadtime 30m

# Create a folder for a longer task (expires in 8 hours)
tempi new --deadtime 8h

# Check what folders you have
tempi list

# Clean up expired folders manually
tempi clean

# Remove everything immediately
tempi deletenow
```

## Technical Details

- **Language**: Go
- **CLI Framework**: Cobra
- **Platform**: Windows
- **Metadata Format**: JSON
- **Process Detection**: Uses Windows API for detecting and terminating processes using folder files

## License

This project is part of a personal utility collection.

