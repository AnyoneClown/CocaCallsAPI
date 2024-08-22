package api

import (
    "encoding/json"
    "net/http"

    "github.com/AnyoneClown/CocaCallsAPI/types"
)

type Response interface {
    Encode() ([]byte, error)
}

type ErrorResponse struct {
    Error string `json:"error"`
    Code  int    `json:"code"`
}


func (r ErrorResponse) Encode() ([]byte, error) {
    return json.Marshal(r)
}

type DefaultResponse struct {
    Message string `json:"message"`
    Code    int    `json:"code"`
}

func (r DefaultResponse) Encode() ([]byte, error) {
    return json.Marshal(r)
}

type UserSuccessResponse struct {
    Message string     `json:"status"`
    User    types.User `json:"user"`
    Code    int        `json:"code"`
}

func (r UserSuccessResponse) Encode() ([]byte, error) {
    return json.Marshal(r)
}

type UserLoginResponse struct {
    Message string                    `json:"status"`
    User    types.UserWithoutPassword `json:"user"`
    Code    int                       `json:"code"`
    Token   string                    `json:"token"`
}

func (r UserLoginResponse) Encode() ([]byte, error) {
    return json.Marshal(r)
}

type UserMEResponse struct {
    Message string                    `json:"status"`
    User    types.UserWithoutPassword `json:"user"`
    Code    int                       `json:"code"`
}

func (r UserMEResponse) Encode() ([]byte, error) {
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

func sendErrorResponse(w http.ResponseWriter, errorMessage string, code int) {
    response := ErrorResponse{
        Error: errorMessage,
        Code:  code,
    }
    sendResponse(w, code, response)
}

func sendSuccessResponse(w http.ResponseWriter, message string, code int) {
    response := DefaultResponse{
        Message: message,
        Code:    code,
    }
    sendResponse(w, code, response)
}

func sendUserSuccessResponse(w http.ResponseWriter, successMessage string, code int, user types.User) {
    response := UserSuccessResponse{
        Message: successMessage,
        User:    user,
        Code:    code,
    }
    sendResponse(w, code, response)
}

func sendUserLoginResponse(w http.ResponseWriter, successMessage string, code int, user types.UserWithoutPassword, token string) {
    response := UserLoginResponse{
        Message: successMessage,
        User:    user,
        Code:    code,
        Token:   token,
    }
    sendResponse(w, code, response)
}

func sendUserMEResponse(w http.ResponseWriter, successMessage string, code int, user types.UserWithoutPassword) {
    response := UserMEResponse{
        Message: successMessage,
        User:    user,
        Code:    code,
    }
    sendResponse(w, code, response)
}