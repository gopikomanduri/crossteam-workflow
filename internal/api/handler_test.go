package api

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTrigonometry(t *testing.T) {
	tests := []struct {
		name      string
		theta     float64
		wantSin   float64
		wantCos   float64
		wantTan   *float64
		wantCosec *float64
		wantSec   *float64
		wantCot   *float64
	}{
		{name: "zero", theta: 0, wantSin: 0, wantCos: 1, wantTan: pointer(0), wantSec: pointer(1)},
		{name: "thirty", theta: 30, wantSin: 0.5, wantCos: math.Sqrt(3) / 2, wantTan: pointer(1 / math.Sqrt(3)), wantCosec: pointer(2), wantSec: pointer(2 / math.Sqrt(3)), wantCot: pointer(math.Sqrt(3))},
		{name: "forty five", theta: 45, wantSin: math.Sqrt(0.5), wantCos: math.Sqrt(0.5), wantTan: pointer(1), wantCosec: pointer(math.Sqrt(2)), wantSec: pointer(math.Sqrt(2)), wantCot: pointer(1)},
		{name: "sixty", theta: 60, wantSin: math.Sqrt(3) / 2, wantCos: 0.5, wantTan: pointer(math.Sqrt(3)), wantCosec: pointer(2 / math.Sqrt(3)), wantSec: pointer(2), wantCot: pointer(1 / math.Sqrt(3))},
		{name: "ninety", theta: 90, wantSin: 1, wantCos: 0, wantCosec: pointer(1), wantCot: pointer(0)},
		{name: "negative", theta: -30, wantSin: -0.5, wantCos: math.Sqrt(3) / 2, wantTan: pointer(-1 / math.Sqrt(3)), wantCosec: pointer(-2), wantSec: pointer(2 / math.Sqrt(3)), wantCot: pointer(-math.Sqrt(3))},
		{name: "over full turn", theta: 390, wantSin: 0.5, wantCos: math.Sqrt(3) / 2, wantTan: pointer(1 / math.Sqrt(3)), wantCosec: pointer(2), wantSec: pointer(2 / math.Sqrt(3)), wantCot: pointer(math.Sqrt(3))},
		{name: "cosec asymptote", theta: 180, wantSin: 0, wantCos: -1, wantTan: pointer(0), wantSec: pointer(-1)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Trigonometry(test.theta)
			assertClose(t, got.Sin, test.wantSin)
			assertClose(t, got.Cos, test.wantCos)
			assertOptionalClose(t, got.Tan, test.wantTan)
			assertOptionalClose(t, got.Cosec, test.wantCosec)
			assertOptionalClose(t, got.Sec, test.wantSec)
			assertOptionalClose(t, got.Cot, test.wantCot)
		})
	}
}

func TestTrigonometryHandler(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		method string
		status int
	}{
		{name: "valid", method: http.MethodPost, body: `{"theta":90}`, status: http.StatusOK},
		{name: "missing theta", method: http.MethodPost, body: `{}`, status: http.StatusBadRequest},
		{name: "null theta", method: http.MethodPost, body: `{"theta":null}`, status: http.StatusBadRequest},
		{name: "unknown field", method: http.MethodPost, body: `{"theta":30,"x":1}`, status: http.StatusBadRequest},
		{name: "invalid JSON", method: http.MethodPost, body: `{`, status: http.StatusBadRequest},
		{name: "wrong method", method: http.MethodGet, status: http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, "/trigonometry", strings.NewReader(test.body))
			TrigonometryHandler(recorder, request)
			if recorder.Code != test.status {
				t.Fatalf("status = %d, want %d", recorder.Code, test.status)
			}
			if test.status == http.StatusOK {
				var response TrigonometryResult
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if response.Tan != nil || response.Sec != nil || response.Cosec == nil || response.Cot == nil {
					t.Fatalf("unexpected asymptote response: %+v", response)
				}
			}
		})
	}
}

func TestCapabilitiesIncludesTrigonometry(t *testing.T) {
	recorder := httptest.NewRecorder()
	Capabilities(recorder, httptest.NewRequest(http.MethodPost, "/capabilities", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response struct {
		Operations []struct {
			ID string `json:"id"`
		} `json:"operations"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	for _, operation := range response.Operations {
		if operation.ID == "trigonometry" {
			return
		}
	}
	t.Fatal("capabilities response does not include trigonometry")
}

func pointer(value float64) *float64 { return &value }

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-12 {
		t.Fatalf("value = %v, want %v", got, want)
	}
}

func assertOptionalClose(t *testing.T, got, want *float64) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("value = %v, want %v", got, want)
		}
		return
	}
	assertClose(t, *got, *want)
}
