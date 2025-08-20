package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/net/html"
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or http)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or http)")
	flag.Parse()

	// Create a new MCP server
	s := server.NewMCPServer(
		"Nix MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
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

	if transport == "http" {
		serv := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP server listening on :8080/mcp")
		go func() {
			if err := serv.Start(":8080"); err != nil {
				log.Printf("Server error: %v", err)
			}
			log.Println("Server execution ended!")
		}()

		var stopChan = make(chan os.Signal, 2)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-stopChan // wait for SIGINT
		log.Println("Shutting down server...")
		err := serv.Shutdown(context.Background())
		if err != nil {
			log.Panic(err)
		}
		log.Println("Server shutdown!")
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

type NixPackageInfo struct {
	Name     string
	Version  string
	HomePage string
	Source   string
}
type SearchPackageResult struct {
	Packages []NixPackageInfo
}

var INTERNAL_MCP_ERROR = errors.New("Internal MCP Error")
var EXTERNAL_ERROR = errors.New("External Error")

func search_package_core(ctx context.Context, name string) (SearchPackageResult, error, bool) {
	result := SearchPackageResult{
		Packages: []NixPackageInfo{},
	}

	// https://search.nixos.org/packages?type=packages&query={name}
	parseUrl, err := url.Parse("https://search.nixos.org/packages?type=packages")
	if err != nil {
		return result, err, true
	}
	queryValues := parseUrl.Query()
	queryValues.Add("query", name)
	parseUrl.RawQuery = queryValues.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", parseUrl.String(), nil)
	if err != nil {
		return result, errors.Join(
			INTERNAL_MCP_ERROR,
			errors.New("Failed to construct request with context"),
			err,
		), true
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, errors.Join(EXTERNAL_ERROR, err), false
	}

	node, err := html.Parse(response.Body)
	if err != nil {
		return result, errors.Join(EXTERNAL_ERROR, err), false
	}

	log.Println(node)

	return result, nil, false
}

func search_package(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, err, should_terminate := search_package_core(ctx, name)
	if err != nil {
		if should_terminate {
			return mcp.NewToolResultError(err.Error()), err
		} else {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
