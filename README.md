[![](https://badge.mcpx.dev?type=server 'MCP Server')](https://github.com/punkpeye/awesome-mcp-servers?tab=readme-ov-file#communication)
[![](https://img.shields.io/badge/OS_Agnostic-Works_Everywhere-purple)](https://github.com/chaindead/tooner?tab=readme-ov-file#installation)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Visitors](https://api.visitorbadge.io/api/visitors?path=https%3A%2F%2Fgithub.com%2Fchaindead%2Ftooner&label=Visitors&labelColor=%23d9e3f0&countColor=%23697689&style=flat&labelStyle=none)](https://visitorbadge.io/status?path=https%3A%2F%2Fgithub.com%2Fchaindead%2Ftooner)

# Tooner

Tooner is a lightweight MCP wrapper that keeps your server setup unchanged, but makes model-facing responses cleaner and shorter. It rewrites JSON-heavy outputs into [TOON format](https://toonformat.dev/) so models spend fewer tokens on syntax noise (~40% fewer tokens).

- [What it does](#what-it-does)
- [Installation](#installation)
  - [Homebrew](#homebrew)
  - [NPX](#npx)
  - [From Releases](#from-releases)
    - [MacOS](#macos)
    - [Linux](#linux)
    - [Windows](#windows)
  - [From Source](#from-source)
- [Configure MCP](#configure-mcp)

## What it does

- Works with your existing MCP servers (no server rewrite required)
- Keeps the same data, but in a format models read more efficiently
- Reduces token usage on large JSON responses
- Handles many messy JSON-like outputs and still returns useful model-ready content

## Installation

### Homebrew

You can install a binary release on macOS/Linux using brew:

```bash
# Install
brew install chaindead/tap/tooner

# Update
brew upgrade chaindead/tap/tooner
```

### NPX

You can run the latest version directly using npx (supports macOS, Linux, and Windows):

```bash
npx -y @chaindead/tooner
```

When using NPX [Configure MCP (Cursor example)](#configure-mcp-cursor-example) becomes:
```json
{
  "mcpServers": {
    "telegram": {
      "command": "npx",
      "args": ["-y", "@chaindead/tooner", "memory-mcp-server-go"]
    }
  }
}
```

### From Releases

#### MacOS

<details>

> **Note:** The commands below install to `/usr/local/bin`. To install elsewhere, replace `/usr/local/bin` with your preferred directory in your PATH.

First, download the archive for your architecture:

```bash
# For Intel Mac (x86_64)
curl -L -o tooner.tar.gz https://github.com/chaindead/tooner/releases/latest/download/tooner_Darwin_x86_64.tar.gz

# For Apple Silicon (M1/M2)
curl -L -o tooner.tar.gz https://github.com/chaindead/tooner/releases/latest/download/tooner_Darwin_arm64.tar.gz
```

Then install the binary:

```bash
# Extract the binary
sudo tar xzf tooner.tar.gz -C /usr/local/bin

# Make it executable
sudo chmod +x /usr/local/bin/tooner

# Clean up
rm tooner.tar.gz
```
</details>

#### Linux
<details>

> **Note:** The commands below install to `/usr/local/bin`. To install elsewhere, replace `/usr/local/bin` with your preferred directory in your PATH.

First, download the archive for your architecture:

```bash
# For x86_64 (64-bit)
curl -L -o tooner.tar.gz https://github.com/chaindead/tooner/releases/latest/download/tooner_Linux_x86_64.tar.gz

# For ARM64
curl -L -o tooner.tar.gz https://github.com/chaindead/tooner/releases/latest/download/tooner_Linux_arm64.tar.gz
```

Then install the binary:

```bash
# Extract the binary
sudo tar xzf tooner.tar.gz -C /usr/local/bin

# Make it executable
sudo chmod +x /usr/local/bin/tooner

# Clean up
rm tooner.tar.gz
```
</details>

#### Windows

<details>

#### Windows
1. Download the latest release for your architecture:
    - [Windows x64](https://github.com/chaindead/tooner/releases/latest/download/tooner_Windows_x86_64.zip)
    - [Windows ARM64](https://github.com/chaindead/tooner/releases/latest/download/tooner_Windows_arm64.zip)
2. Extract the `.zip` file
3. Add the extracted directory to your PATH or move `tooner.exe` to a directory in your PATH
</details>

### From Source

Requirements:
- Go 1.26 or later
- GOBIN in PATH

```bash
go install github.com/chaindead/tooner@latest
```

## Configure MCP

Setup is the only thing you need to do: in MCP config, shift the command by one argument and insert `tooner` at the beginning.

Before:

```json
{
  "mcpServers": {
    "memory": {
      "command": "memory-mcp-server-go"
    }
  }
}
```

After:

```json
{
  "mcpServers": {
    "memory": {
      "command": "tooner",
      "args": ["memory-mcp-server-go"]
    }
  }
}
```

Another example - Postgres MCP:

```json
{
  "mcpServers": {
    "postgres": {
      "command": "tooner",
      "args": [
        "uvx",
        "postgres-mcp",
        "--access-mode=unrestricted"
      ],
      "env": {
        "DATABASE_URI": "...",
        "TOONER_LOG_PATH": "/tmp/tooner.log"
      }
    }
  }
}
```

NPX setup:

```json
{
  "mcpServers": {
    "memory": {
      "command": "npx",
      "args": ["-y", "@chaindead/tooner", "memory-mcp-server-go"]
    }
  }
}
```

That's it. From there, `tooner` launches your MCP server as a subprocess and transparently proxies traffic, replacing JSON with TOON.
