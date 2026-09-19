package api

import (
	"encoding/json"
	"net/http"
)

type operands struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type resultResp struct {
	Result float64 `json:"result"`
}

type errorResp struct {
	Error string `json:"error"`
}

type trigonometryRequest struct {
	Theta *float64 `json:"theta" binding:"required"`
}

func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Capabilities godoc
// @Summary List supported calculator operations
// @Description Returns the calculator operations supported by this service.
// @Tags calculator
// @Produce json
// @Success 200 {object} map[string][]map[string]string
// @Failure 405 {object} errorResp
// @Router /capabilities [post]
func Capabilities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResp{Error: "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	ops := []map[string]string{
		{"id": "add", "description": "Add two numbers"},
		{"id": "subtract", "description": "Subtract second number from first"},
		{"id": "multiply", "description": "Multiply two numbers"},
		{"id": "divide", "description": "Divide first number by second"},
		{"id": "trigonometry", "description": "Evaluate trigonometric functions for an angle in degrees"},
	}
	writeJSON(w, map[string]interface{}{"operations": ops}, http.StatusOK)
}

// TrigonometryHandler godoc
// @Summary Evaluate trigonometric functions for an angle in degrees
// @Description Returns sin, cos, tan, cosec, sec, and cot. Undefined reciprocal values are null.
// @Tags calculator
// @Accept json
// @Produce json
// @Param request body trigonometryRequest true "Angle in degrees"
// @Success 200 {object} TrigonometryResult
// @Failure 400 {object} errorResp
// @Failure 405 {object} errorResp
// @Router /trigonometry [post]
func TrigonometryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResp{Error: "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	var request trigonometryRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&request); err != nil || request.Theta == nil {
		writeJSON(w, errorResp{Error: "invalid JSON body"}, http.StatusBadRequest)
		return
	}
	writeJSON(w, Trigonometry(*request.Theta), http.StatusOK)
}

func decodeOperands(r *http.Request) (operands, *errorResp) {
	var o operands
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&o); err != nil {
		return o, &errorResp{Error: "invalid JSON body"}
	}
	return o, nil
}

// Add godoc
// @Summary Add two numbers
// @Tags calculator
// @Accept json
// @Produce json
// @Param operands body operands true "Operands"
// @Success 200 {object} resultResp
// @Failure 400 {object} errorResp
// @Failure 405 {object} errorResp
// @Router /add [post]
func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResp{Error: "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	o, err := decodeOperands(r)
	if err != nil {
		writeJSON(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, resultResp{Result: o.X + o.Y}, http.StatusOK)
}

// Subtract godoc
// @Summary Subtract the second number from the first
// @Tags calculator
// @Accept json
// @Produce json
// @Param operands body operands true "Operands"
// @Success 200 {object} resultResp
// @Failure 400 {object} errorResp
// @Failure 405 {object} errorResp
// @Router /subtract [post]
func Subtract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResp{Error: "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	o, err := decodeOperands(r)
	if err != nil {
		writeJSON(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, resultResp{Result: o.X - o.Y}, http.StatusOK)
}

// Multiply godoc
// @Summary Multiply two numbers
// @Tags calculator
// @Accept json
// @Produce json
// @Param operands body operands true "Operands"
// @Success 200 {object} resultResp
// @Failure 400 {object} errorResp
// @Failure 405 {object} errorResp
// @Router /multiply [post]
func Multiply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResp{Error: "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	o, err := decodeOperands(r)
	if err != nil {
		writeJSON(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, resultResp{Result: o.X * o.Y}, http.StatusOK)
}

// Divide godoc
// @Summary Divide the first number by the second
// @Tags calculator
// @Accept json
// @Produce json
// @Param operands body operands true "Operands"
// @Success 200 {object} resultResp
// @Failure 400 {object} errorResp
// @Failure 405 {object} errorResp
// @Router /divide [post]
func Divide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResp{Error: "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	o, err := decodeOperands(r)
	if err != nil {
		writeJSON(w, err, http.StatusBadRequest)
		return
	}
	if o.Y == 0 {
		writeJSON(w, errorResp{Error: "division by zero"}, http.StatusBadRequest)
		return
	}
	writeJSON(w, resultResp{Result: o.X / o.Y}, http.StatusOK)
}
