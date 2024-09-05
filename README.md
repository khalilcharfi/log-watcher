# Log File Watcher

A Go-based utility for watching log files and directories for changes. The tool supports filtering log entries, prettifying output with color-coding, and handling log files in JSON format.

![go](https://img.shields.io/static/v1?label=Golang+1.23&labelColor=34a1eb&message=Go&color=000000&logo=go&logoColor=ffffff&style=flat-square) 
![go](https://img.shields.io/static/v1?label=Channels&labelColor=34a1eb&message=Go&color=000000&logo=go&logoColor=ffffff&style=flat-square)

<p align="center">
<img src="https://raw.githubusercontent.com/khalilcharfi/log-watcher/assets/logo.webp" height="300">
</p>
<br>

## Features

- Watch individual log files or entire directories for changes.
- Filter logs by specific key-value pairs (e.g., `log_level[debug]`).
- Color-code log entries for easy differentiation.
- Handles logs in JSON format and extracts key details for display.

## Requirements

- **Go 1.23**

## How to Install it:

First, clone this repository to your machine:

~~~bash
git clone https://github.com/khalilcharfi/log-watcher.git
~~~

Next, run the following command to build and install the application:

~~~bash
chmod +x install.sh
./install.sh
~~~

This script will:
- Check if Go is installed and ensure it’s version 1.20 or higher.
- Download and install Go if it’s not already installed.
- Build the `log-watcher` binary for your system.
- Optionally, add an alias for `log-watcher` to your bash or zsh shell.

## How to Use it:

### In Development Mode:

When developing, you can run the tool directly using `go run`:

~~~bash
go run ./cmd {filepath}
~~~

Replace `{filepath}` with the path to the log file you want to watch.

### After Compiling:

After compiling the application, you can run it as follows:

~~~bash
./log-watcher {filepath}
~~~

Replace `{filepath}` with the path to the log file you want to watch.

## Advanced Usage

### Watch a Directory for Log Files

To watch a directory for `.log` files, use the `-d` flag:

~~~bash
./log-watcher -d /path/to/your/log/directory
~~~

This will watch all `.log` files in the specified directory. If any new `.log` files are created, they will automatically be watched.

### Prettify Output with Color Coding

To enable color-coded output for each file, use the `-p` flag:

~~~bash
./log-watcher -p /path/to/your/log/file.log
~~~

### Filter Logs by Specific Key-Value Pairs

You can filter logs by a specific key-value pair using the `--filter` flag. For example:

~~~bash
./log-watcher --filter=log_level[debug] /path/to/your/log/file.log
~~~

This command will only display logs where `log_level` is `debug`.

### Combining Options

All options can be combined. For example:

~~~bash
./log-watcher -d /path/to/your/log/directory -p --filter=category[debug]
~~~

This will watch all `.log` files in the specified directory, filter logs by `category=debug`, and apply color-coded output.

### Default Behavior

If you simply pass a file path as an argument (e.g., `go run ./cmd {filepath}` or `./log-watcher {filepath}`), the program will watch that specific file without any prettification or filtering.
