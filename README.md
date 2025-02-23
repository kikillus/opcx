# OPCX

A Terminal User Interface (TUI) explorer for OPC UA servers, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [gopcua](https://github.com/gopcua/opcua).

## Overview

OPCX provides an interactive terminal interface for browsing and exploring OPC UA servers. This project is currently heavily work in progress.

## Features

- Terminal-based UI using Bubble Tea framework
- Browse OPC UA server nodes
- Recursive node exploration
- Read node values
- Interactive navigation
- Monitor node values

## Requirements

- Go 1.24 or higher
- Dependencies:
  - github.com/charmbracelet/bubbletea
  - github.com/charmbracelet/bubbles
  - github.com/gopcua/opcua

## Installation

```bash
go install github.com/yourusername/opcx/cmd/opcx@latest
```

## Usage

```bash
opcx
```

## Development Status

This project is under active development and should be considered not even alpha quality. Features and API may change without notice.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.