package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AnyoneClown/CocaCallsAPI/types"
	"github.com/AnyoneClown/CocaCallsAPI/utils"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserByIDRequest struct {
	UserID string `json:"user_id"`
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

	sendUserSuccessResponse(w, "Successfully created a user", http.StatusCreated, user)
}

func (s *Server) handleUserMe(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		sendErrorResponse(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		sendErrorResponse(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	claims, err := utils.ExtractClaimsFromToken(tokenString)
	if err != nil {
		sendErrorResponse(w, "Can't extract claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		sendErrorResponse(w, "Email not found in token claims", http.StatusUnauthorized)
		return
	}

	user, err := s.storage.GetUserByEmail(email)
	if err != nil {
		sendErrorResponse(w, "Can't find user with this email", http.StatusBadRequest)
		return
	}

	userResponse := types.UserWithoutPassword{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	sendUserMEResponse(w, "User data retrieved successfully", http.StatusOK, userResponse)
}

func (s *Server) handleMainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Main page")
}

func (s *Server) handleJWTCreate(w http.ResponseWriter, r *http.Request) {
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

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		sendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	userResponse := types.UserWithoutPassword{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	sendUserLoginResponse(w, "Successfully logged in", http.StatusOK, userResponse, token)
}

func (s *Server) handleJWTVerify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		sendErrorResponse(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		sendErrorResponse(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	_, err := utils.VerifyToken(tokenString)
	if err != nil {
		sendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	sendSuccessResponse(w, "Token is valid", http.StatusOK)
}

func (s *Server) signInWithProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	log.Printf("Sign in with provider: %s", provider)

	gothic.BeginAuthHandler(w, r)
}

func (s *Server) callbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	log.Printf("Callback for provider: %s", provider)

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		sendErrorResponse(w, "Failed to complete authentication", http.StatusInternalServerError)
		return
	}
	log.Printf("User: %+v", user)
	sendSuccessResponse(w, "User logged successfully", http.StatusOK)
}
