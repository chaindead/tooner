# Tooner

An MCP (Model Context Protocol) proxy that wraps MCP servers and converts JSON responses to [TOON format](https://toonformat.dev/) — a token-efficient alternative to JSON optimized for LLMs (~40% fewer tokens).

- [What it does](#what-it-does)
- [Installation](#installation)
  - [Homebrew](#homebrew)
  - [NPX](#npx)
  - [From Releases](#from-releases)
    - [MacOS](#macos)
    - [Linux](#linux)
    - [Windows](#windows)
  - [From Source](#from-source)
- [Configure MCP (Cursor example)](#configure-mcp-cursor-example)

## What it does

Tooner runs any MCP server as a subprocess and transparently proxies messages between the client (e.g. Cursor) and the server. When the server returns JSON in `tools/call` responses, Tooner converts that JSON to TOON before forwarding it to the client, reducing token usage while preserving the same data.

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

When using NPX, modify the standard commands and configuration as follows:

- [Authentication command](#authorization) becomes:
```bash
npx -y @chaindead/tooner auth ...
```

- [Claude MCP server configuration](#client-configuration) becomes:
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

## Configure MCP (Cursor example)

Add Tooner as a wrapper in `~/.cursor/mcp.json`, passing your MCP server as the first argument:

```json
{
  "mcpServers": {
    "memory": {
      "command": "tooner",
      "args": ["memory-mcp-server-go"],
      "env": {
        "TOONER_LOG_PATH": "/tmp/mcp.log"
      }
    }
  }
}
```

- Replace `memory-mcp-server-go` with any MCP server binary in your PATH (e.g. `go-mcp-postgres`).
- You can pass args and envs to MCP as always
- TOONER_LOG_PATH is optional
