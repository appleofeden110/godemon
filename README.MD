# GoDemon

A file system watcher written in Go that monitors your project directory for changes and reports file modifications, additions, and deletions in real-time.

## Overview

GoDemon is an active development project that traverses a project's directory tree using a **General Tree** data structure and **Depth-First Search (DFS)** algorithm to track file system changes. It watches for modifications and provides color-coded output to notify you of what's happening in your project.

**Status**: Under active development - features are being refined and expanded.

## Features

- **Directory Tree Traversal**: Uses a General Tree structure with DFS to efficiently scan your entire project
- **Real-time File Monitoring**: Continuously watches for file changes every 2 seconds
- **Project Size Tracking**: Displays the total number of files being monitored
- **Configurable Ignore Patterns**: Uses JSON configuration files to exclude directories from monitoring
- **Scalable Architecture**: Built with a queue-based system for efficient file processing

## How It Works

1. **Initialization** (`init.go`): Sets up the file tracking system and begins monitoring
2. **Tree Building** (`tree.go`): Recursively builds a file tree structure of your project directory
3. **Change Detection** (`utils.go`): Compares file states between polling intervals to detect:
   - New or modified files
   - Deleted files and directories
4. **Reporting**: Outputs formatted change notifications with file paths and status

## Project Structure

```
godemon/
├── main.go              # Entry point
├── go.mod              # Go module definition (Go 1.22.2+)
├── tree_check/         # Core monitoring package
│   ├── init.go         # TreeCheck() main function
│   ├── tree.go         # File tree node structures and DFS logic
│   └── utils.go        # Helper functions and file comparison logic
├── ignored.json        # Configuration for directories to ignore
├── godemon             # Compiled binary (git-ignored)
└── README.MD           # This file
```

## Configuration

Create an `ignored.json` file in your project root to specify directories to skip:

```json
{
  "node_modules": true,
  ".git": true,
  ".vscode": true,
  ".idea": true,
  "dist": true,
  "build": true
}
```

## Usage

### Prerequisites
- Go 1.22.2 or higher

### Running GoDemon

```bash
go run main.go
```

Or build and run the binary:

```bash
go build -o godemon
./godemon
```

## Technical Details

- **Data Structure**: General Tree with node-based hierarchy
- **Algorithm**: Depth-First Search (DFS) for directory traversal
- **Polling Interval**: 2 seconds (configurable in code)
- **Configuration Format**: JSON for ignore patterns
- **Output**: Color-formatted console output with ANSI escape codes

## Contributing

Contributions are welcome! Feel free to:
- Report issues or bugs
- Suggest improvements and new features
- Submit pull requests with enhancements

## License

This project is licensed under the MIT License. See [LICENSE.txt](LICENSE.txt) for details.

## Contact

Have questions or suggestions? Reach out:
- **Email**: nazarkaniuka6@gmail.com

---
