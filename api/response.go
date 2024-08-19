package api

import (
	"encoding/json"
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/types"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
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

func sendUserCreatedResponse(w http.ResponseWriter, successMessage string, code int, user types.User) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := struct {
		Message string     `json:"status"`
		User    types.User `json:"user"`
		Code    int        `json:"code"`
	}{
		Message: successMessage,
		User:    user,
		Code:    code,
	}
	json.NewEncoder(w).Encode(response)
}
