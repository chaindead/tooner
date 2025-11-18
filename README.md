# Tooner

An MCP (Model Context Protocol) proxy that wraps MCP servers and converts JSON responses to [TOON format](https://toonformat.dev/) — a token-efficient alternative to JSON optimized for LLMs (~40% fewer tokens).

## What it does

Tooner runs any MCP server as a subprocess and transparently proxies messages between the client (e.g. Cursor) and the server. When the server returns JSON in `tools/call` responses, Tooner converts that JSON to TOON before forwarding it to the client, reducing token usage while preserving the same data.

### Configure Cursor MCP

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

Replace `memory-mcp-server-go` with any MCP server binary in your PATH (e.g. `go-mcp-postgres`, `pocketbase-cursor-mcp`).
