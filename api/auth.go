package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
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
		sendErrorResponse(w, "User with this email already exists", http.StatusBadRequest)
		return
	}

	sendUserCreatedResponse(w, "Successfully created a user", http.StatusCreated, user)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := s.storage.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		log.Printf("Password does not match for user: %v", req.Email)
		sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		sendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	sendUserLoginResponse(w, "Successfully logged in", http.StatusOK, user, token)
}
