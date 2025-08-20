package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "http", "Transport type (stdio or http)")
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
		serverCtx, cancelServerCtx := context.WithCancel(context.Background())
		defer cancelServerCtx()
		serv := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP server listening on :8080/mcp")

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := serv.Start(":8080"); err != nil {
				log.Printf("Server error: %v", err)
			}
			log.Println("Server execution ended!")
		}()

		var stopChan = make(chan os.Signal, 2)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-stopChan // wait for SIGINT
		log.Println("Shutting down server...")
		err := serv.Shutdown(serverCtx)
		if err != nil {
			log.Panic(err)
		}
		wg.Wait()
		log.Println("Server shutdown!")
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

type SearchPackageResult struct {
	Packages []NixHubPkgSummary
}

func search_package(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err, should_terminate := search_package_core(ctx, name)
	if err != nil {
		if should_terminate {
			return mcp.NewToolResultError(err.Error()), err
		} else {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	b := strings.Builder{}
	for _, v := range result.Packages {
		b.WriteString(v.Name)
		b.WriteString(v.Summary)
		b.WriteString("; Last Updated: ")
		b.WriteString(v.LastUpdated.String())
		b.WriteRune('\n')
	}

	return mcp.NewToolResultText(b.String()), nil
}
