package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		x         float64
		y         float64
		want      float64
		wantErr   string
	}{
		{name: "add", operation: "add", x: 3.5, y: 2, want: 5.5},
		{name: "subtract", operation: "subtract", x: 10, y: 4, want: 6},
		{name: "multiply", operation: "multiply", x: 3, y: 4, want: 12},
		{name: "divide", operation: "divide", x: 10, y: 4, want: 2.5},
		{name: "division by zero", operation: "divide", x: 10, y: 0, wantErr: "division by zero"},
		{name: "unknown operation", operation: "modulo", x: 10, y: 4, wantErr: "unsupported operation: modulo"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := calculate(test.operation, test.x, test.y)
			if test.wantErr != "" {
				if err == nil || err.Error() != test.wantErr {
					t.Fatalf("calculate() error = %v, want %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("calculate() unexpected error: %v", err)
			}
			if got != test.want {
				t.Errorf("calculate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestTrigonometryToolHandler(t *testing.T) {
	result, err := trigonometryHandler(context.Background(), mcp.CallToolRequest{}, trigonometryArgs{Theta: 90})
	if err != nil {
		t.Fatalf("trigonometryHandler() error: %v", err)
	}
	if result.IsError {
		t.Fatalf("trigonometryHandler() returned an MCP tool error")
	}

	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
	if payload["tan"] != nil || payload["sec"] != nil {
		t.Fatalf("undefined values must be null, got %s", encoded)
	}
	if payload["sin"] != float64(1) || payload["cos"] != float64(0) {
		t.Fatalf("unexpected cardinal-angle values: %s", encoded)
	}
}
