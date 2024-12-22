package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result *float64 `json:"result,omitempty"`
	Error  *string  `json:"error,omitempty"`
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result, err := Calc(req.Expression)
	response := Response{}

	if err != nil {
		if errors.Is(err, ErrInvalidExpression) || errors.Is(err, ErrMismatchedParentheses) {
			msg := "Expression is not valid"
			response.Error = &msg
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(response)
			return
		} else {
			msg := "Internal server error"
			response.Error = &msg
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
	}

	response.Result = &result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v1/calculate", handleCalculate)
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
