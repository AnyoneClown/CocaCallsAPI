package utils

import (
	"encoding/json"
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/types"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (r Response) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func sendResponse(w http.ResponseWriter, code int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encodedResponse, err := response.Encode()
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Write(encodedResponse)
}

func SendErrorResponse(w http.ResponseWriter, errorMessage string, code int) {
	response := Response{
		Code:  code,
		Error: errorMessage,
	}
	sendResponse(w, code, response)
}

func SendSuccessResponse(w http.ResponseWriter, message string, code int) {
	response := Response{
		Code:    code,
		Message: message,
	}
	sendResponse(w, code, response)
}

func SendDataResponse(w http.ResponseWriter, message string, code int, data interface{}) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	sendResponse(w, code, response)
}

func SendUserResponse(w http.ResponseWriter, message string, code int, user types.UserWithoutPassword) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    user,
	}
	sendResponse(w, code, response)
}

func SendUserLoginResponse(w http.ResponseWriter, message string, code int, user types.UserWithoutPassword, token string) {
	response := Response{
		Code:    code,
		Message: message,
		Data: struct {
			User  types.UserWithoutPassword `json:"user"`
			Token string                    `json:"token"`
		}{
			User:  user,
			Token: token,
		},
	}
	sendResponse(w, code, response)
}
