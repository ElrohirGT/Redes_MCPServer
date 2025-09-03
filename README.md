# Nix Package MCP Server

https://github.com/user-attachments/assets/1b838292-ab21-4943-a3fb-8fd5f007ddb9

## How to use?

The MCP server was built using the Go programming language!

**ENDPOINT:** The MCP server default endpoint for communication is `/mcp`.

If you have Nix installed running this server is as simple as:

```bash
# If you're inside the root of the repo
nix run
```

```bash
# If you don't have the repo cloned
nix run github:ElrohirGT/Redes_MCPServer
```

You can also just compile a binary using normal go commands:

```bash
# If you want to just build a binary
go build .
```

```bash
# If you want to just run the program
go run . 
```

This MCP server supports being run with HTTP or with STDIN/STDOUT, please run
the following command to check the available options:

```bash
go run . -h
```

## How to change the server address

By default the server runs on: `0.0.0.0:8080` but you can change this using the
following flag:

```bash
go run . -a <custom host and port>
```

Again you can check all available configuration options by using:

```bash
go run . -h
```

## Available Tools

This MCP Server exposes a single tool, this tool allows to connected clients to
search for packages inside the Nix repository, and obtain info like descriptions
and the date when the package was last updated.
