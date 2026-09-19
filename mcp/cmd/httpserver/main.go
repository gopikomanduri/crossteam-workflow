package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"samplecalculatorproject/internal/api"
)

type calculateArgs struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type result struct {
	Result float64 `json:"result"`
}

type trigonometryArgs struct {
	Theta float64 `json:"theta"`
}

func main() {
	mcpServer := server.NewMCPServer(
		"sample-calculator",
		"1.1.0",
		server.WithToolCapabilities(false),
	)

	addCalculatorTools(mcpServer)

	addr := envOrDefault("MCP_HTTP_ADDR", ":8081")
	httpServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithEndpointPath("/mcp"),
		server.WithStateLess(true),
	)

	log.Printf("starting MCP HTTP server on http://localhost%s/mcp", addr)
	if err := httpServer.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("MCP HTTP server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func addCalculatorTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(
		mcp.NewTool("trigonometry",
			mcp.WithDescription("Evaluate sin, cos, tan, cosec, sec, and cot for an angle in degrees. Undefined reciprocal values are null."),
			mcp.WithNumber("theta", mcp.Required(), mcp.Description("Angle in degrees.")),
		),
		mcp.NewTypedToolHandler(trigonometryHandler),
	)
	mcpServer.AddTool(
		mcp.NewTool("add",
			mcp.WithDescription("Add two numbers."),
			mcp.WithNumber("x", mcp.Required(), mcp.Description("First number.")),
			mcp.WithNumber("y", mcp.Required(), mcp.Description("Second number.")),
		),
		mcp.NewTypedToolHandler(operationHandler("add")),
	)
	mcpServer.AddTool(
		mcp.NewTool("subtract",
			mcp.WithDescription("Subtract the second number from the first."),
			mcp.WithNumber("x", mcp.Required(), mcp.Description("First number.")),
			mcp.WithNumber("y", mcp.Required(), mcp.Description("Second number.")),
		),
		mcp.NewTypedToolHandler(operationHandler("subtract")),
	)
	mcpServer.AddTool(
		mcp.NewTool("multiply",
			mcp.WithDescription("Multiply two numbers."),
			mcp.WithNumber("x", mcp.Required(), mcp.Description("First number.")),
			mcp.WithNumber("y", mcp.Required(), mcp.Description("Second number.")),
		),
		mcp.NewTypedToolHandler(operationHandler("multiply")),
	)
	mcpServer.AddTool(
		mcp.NewTool("divide",
			mcp.WithDescription("Divide the first number by the second."),
			mcp.WithNumber("x", mcp.Required(), mcp.Description("First number.")),
			mcp.WithNumber("y", mcp.Required(), mcp.Description("Second number.")),
		),
		mcp.NewTypedToolHandler(operationHandler("divide")),
	)
	mcpServer.AddTool(
		mcp.NewTool("capabilities",
			mcp.WithDescription("List supported calculator operations."),
		),
		capabilitiesHandler,
	)
}

func trigonometryHandler(_ context.Context, _ mcp.CallToolRequest, args trigonometryArgs) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultJSON(api.Trigonometry(args.Theta))
}

func operationHandler(operation string) func(context.Context, mcp.CallToolRequest, calculateArgs) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest, args calculateArgs) (*mcp.CallToolResult, error) {
		value, err := calculate(operation, args.X, args.Y)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultJSON(result{Result: value})
	}
}

func capabilitiesHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultJSON(map[string][]map[string]string{
		"operations": {
			{"id": "add", "description": "Add two numbers"},
			{"id": "subtract", "description": "Subtract second number from first"},
			{"id": "multiply", "description": "Multiply two numbers"},
			{"id": "divide", "description": "Divide first number by second"},
			{"id": "trigonometry", "description": "Evaluate trigonometric functions for an angle in degrees"},
		},
	})
}

func calculate(operation string, x, y float64) (float64, error) {
	switch operation {
	case "add":
		return x + y, nil
	case "subtract":
		return x - y, nil
	case "multiply":
		return x * y, nil
	case "divide":
		if y == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return x / y, nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", operation)
	}
}
