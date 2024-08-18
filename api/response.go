package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

/*
HTTP Response handling for errors,
Returns valid JSON with error message and response code
*/

func sendErrorResponse(w http.ResponseWriter, errorMessage string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    response := ErrorResponse{
        Error: errorMessage,
        Code:  code,
    }
    json.NewEncoder(w).Encode(response)
}
