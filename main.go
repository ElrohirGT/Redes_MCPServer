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
	var addr string
	flag.StringVar(&addr, "a", ":8080", "Address where the server will listen")
	flag.Parse()

	// Create a new MCP server
	s := server.NewMCPServer(
		"Nix MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool("search_package",
		mcp.WithDescription("Search for packages inside the nix repository, get information like: name, summary, home page url, version, license, release date and platforms"),
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
		serv := server.NewStreamableHTTPServer(
			s,
			server.WithStateLess(true),
		)
		log.Println("HTTP server listening on", addr)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := serv.Start(addr); err != nil {
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
	Packages []NixHubPkgInfo
}

func search_package(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request to `search_package`\n%#v", request)
	name, err := request.RequireString("name")
	if err != nil {
		log.Println("ERROR: No required `name` argument provided!")
		return mcp.NewToolResultError("ERROR: " + err.Error()), nil
	}

	log.Println("Trying to search packages with name:", name)
	result, err, should_terminate := search_package_core(ctx, name)
	if err != nil {
		if should_terminate {
			return mcp.NewToolResultError(err.Error()), err
		} else {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	log.Println("Transforming request from JSON into single string")
	b := strings.Builder{}
	for _, v := range result.Packages {
		b.WriteString(v.Name)
		b.WriteString("\n - Summary: ")
		b.WriteString(v.Summary)
		b.WriteString("\n - Homepage: ")
		b.WriteString(v.HomepageUrl)
		b.WriteString("\n - License: ")
		b.WriteString(v.License)

		if len(v.Releases) >= 1 {
			lr := v.Releases[0]
			b.WriteString("\n - Latest Release Version: ")
			b.WriteString(lr.Version)
			b.WriteString("\n - Latest Release Platforms: ")
			b.WriteString(lr.PlatformsSummary)
			b.WriteString("\n - Latest Release Date: ")
			b.WriteString(lr.LastUpdated.String())
		}
		b.WriteRune('\n')
	}

	resultText := b.String()
	log.Println("Final output:\n", resultText)

	return mcp.NewToolResultText(resultText), nil
}
