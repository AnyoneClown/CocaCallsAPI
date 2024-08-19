package api

import (
	"encoding/json"
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		sendErrorResponse(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := s.storage.CreateUser(req.Email, hashedPassword)
	if err != nil {
		sendErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
	}

	sendUserCreatedResponse(w, "Successfully created a user", http.StatusCreated, user)
}
