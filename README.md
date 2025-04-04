# MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/luffy050596/mcp-server)](https://goreportcard.com/report/github.com/luffy050596/mcp-server)
[![GoDoc](https://godoc.org/github.com/luffy050596/mcp-server?status.svg)](https://godoc.org/github.com/luffy050596/mcp-server)
[![License](https://img.shields.io/github/license/luffy050596/mcp-server.svg)](https://github.com/luffy050596/mcp-server/blob/main/LICENSE)

MCP Server is a Go-based MCP tools project that provides a series of MCP utilities. 
This is a personal learning project, please do not use it in production environment.
Using [github.com/ThinkInAIXYZ/go-mcp](https://github.com/ThinkInAIXYZ/go-mcp) as the MCP core framework.

## Features

- IP address processing service
- Time service
- Poster generation service
- More features coming soon...

## System Requirements

- Go 1.23.0 or higher
- Supports Linux, macOS and Windows

## Quick Start

### Installation

```bash
# Clone the project
git clone https://github.com/luffy050596/mcp-server.git
cd mcp-server

# Install dependencies
go mod download
```

### Build

Build a single service:
```bash
make build dir=<service_directory>
```

Build all services:
```bash
make build-all
```

### Test

Run all tests:
```bash
make test
```

### Run

Run a single service:
```bash
./bin/mcp-ip -mode=stdio -addr=:59001
```

#### Parameters

- `-mode` Running mode, available values are `stdio` or `sse`. Default is `stdio`
- `-addr` Service address, required when `-mode=sse`
- `-key` Bailian API Key, required for poster service

## Project Structure

```
.
├── bin/           # Compiled binary files
├── ip/            # IP geolocation info, using https://ip.rpcx.io API
├── time/          # Time query and timestamp conversion
├── poster/        # Poster generation service, using Bailian Creative Poster Generation API(https://help.aliyun.com/zh/model-studio/creative-poster-generation)
└── pkg/           # Shared packages
```

## Development Tools

The project uses the following development tools to ensure code quality:

- `.golangci.yaml` - golangci-lint configuration
- `.pre-commit-config.yaml` - Git pre-commit hooks
- `.gitleaks.toml` - Gitleaks sensitive information detection configuration

## Dependency Management

Main dependencies:

- github.com/ThinkInAIXYZ/go-mcp - MCP core library

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the terms of the LICENSE file - see the [LICENSE](LICENSE) file for details.

## Contact

- Project Maintainer: [Your Name]
- Project Link: [https://github.com/luffy050596/mcp-server](https://github.com/luffy050596/mcp-server)

## Acknowledgments

Thanks to all developers who have contributed to this project.
