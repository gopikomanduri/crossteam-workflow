package main

import "testing"

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
