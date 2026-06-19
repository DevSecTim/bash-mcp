# Bash MCP Server

A lightweight [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server written in Go that allows LLM agents to execute bash commands on the host machine. 

Since MCP natively communicates over standard input/output (`stdio`), this server can be run entirely locally or transparently over SSH without any extra proxy configurations.

## Building

To build the server, ensure you have Go installed, then run:

```bash
make build
```

This will produce an executable binary at `bin/bash-mcp`.

## Usage

You can use this server with any MCP-compatible client. The client will spawn the binary and communicate with it using JSON-RPC messages over `stdin` and `stdout`.

### 1. Local Execution

To run the server locally on the same machine as your MCP client, configure your client to execute the compiled binary directly.

**Example `mcp.json` configuration:**
```json
{
  "mcpServers": {
    "bash-local": {
      "command": "/absolute/path/to/bash-mcp/bin/bash-mcp",
      "args": []
    }
  }
}
```

### 2. Remote Execution via SSH

Because SSH automatically proxies `stdin` and `stdout`, you can run this server on a remote machine over SSH natively. The client will simply launch the `ssh` command, and the remote execution will seamlessly connect the IO streams.

*Requirements:* Make sure you have passwordless SSH authentication (e.g., SSH keys) set up between your client machine and the remote host, as the MCP client cannot interactively enter passwords.

**Example `mcp.json` configuration:**
```json
{
  "mcpServers": {
    "bash-remote": {
      "command": "ssh",
      "args": [
        "user@remote-host",
        "/absolute/path/to/bash-mcp/bin/bash-mcp"
      ]
    }
  }
}
```

### Tools Provided

- `run_bash_command`: Executes a bash command and returns its combined `stdout` and `stderr`.
