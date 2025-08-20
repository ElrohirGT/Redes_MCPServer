package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/net/html"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Nix MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool("search_package",
		mcp.WithDescription("Search for a specific package inside Nix repository"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the package to search"),
		),
	)

	// Add tool handler
	s.AddTool(tool, search_package)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func search_package(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	// https://search.nixos.org/packages?type=packages&query={name}
	parseUrl, err := url.Parse("https://search.nixos.org/packages?type=packages")
	if err != nil {
		return mcp.NewToolResultError("Failed to parse URL: " + err.Error()), err
	}
	queryValues := parseUrl.Query()
	queryValues.Add("query", name)
	parseUrl.RawQuery = queryValues.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", parseUrl.String(), nil)
	if err != nil {
		return mcp.NewToolResultError("Failed to construct request with context: " + err.Error()), err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	node, err := html.Parse(response.Body)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	log.Println(node)

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
