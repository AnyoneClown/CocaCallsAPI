package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AnyoneClown/CocaCallsAPI/types"
	"github.com/AnyoneClown/CocaCallsAPI/utils"
	"golang.org/x/oauth2"
)

type AuthRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	GoogleID      string `json:"google_id,omitempty"`
	Picture       string `json:"picture,omitempty"`
	Provider      string `json:"provider,omitempty"`
	VerifiedEmail bool   `json:"verified_email,omitempty"`
}

func (h *DefaultHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	userBody := &types.UserToCreate{
		Email:         req.Email,
		Password:      hashedPassword,
		GoogleID:      req.GoogleID,
		Picture:       req.Picture,
		Provider:      req.Provider,
		VerifiedEmail: req.VerifiedEmail,
		IsAdmin:       false,
	}
	user, err := h.Storage.CreateUser(*userBody)
	if err != nil {
		utils.SendErrorResponse(w, "User with this email already exists", http.StatusBadRequest)
		return
	}

	userResponse := types.UserWithoutPassword{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.SendUserResponse(w, "Successfully created a user", http.StatusCreated, userResponse)
}

func (h *DefaultHandler) HandleJWTCreate(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.Storage.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		utils.SendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		log.Printf("Password does not match for user: %v", req.Email)
		utils.SendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		utils.SendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	userResponse := types.UserWithoutPassword{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.SendUserLoginResponse(w, "Successfully logged in", http.StatusOK, userResponse, token)
}

func (h *DefaultHandler) HandleJWTVerify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.SendErrorResponse(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		utils.SendErrorResponse(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	_, err := utils.VerifyToken(tokenString)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	utils.SendSuccessResponse(w, "Token is valid", http.StatusOK)
}

func (h *DefaultHandler) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := utils.GenerateStateOauthCookie(w)
	u := utils.GoogleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (h *DefaultHandler) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.URL.Query().Get("error"); err != "" {
		frontendURL := utils.GetEnvVariable("FRONTEND_URL")
		http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
		return
	}

	data, err := utils.GetUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var userInfo types.GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		log.Println("Failed to parse user info:", err)
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	user, err := h.Storage.GetUserByEmail(userInfo.Email)
	if err != nil {
		// User does not exist, create a new user
		userBody := &types.UserToCreate{
			Email:         userInfo.Email,
			Password:      "",
			GoogleID:      userInfo.GoogleID,
			Picture:       userInfo.Picture,
			Provider:      "google",
			VerifiedEmail: userInfo.VerifiedEmail,
		}
		user, err = h.Storage.CreateUser(*userBody)
		if err != nil {
			log.Println("Failed to create user:", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if user.GoogleID == "" {
		// User exists but does not have GoogleID, update the user
		user.GoogleID = userInfo.GoogleID
		user.Picture = userInfo.Picture
		user.Provider = "google"
		user.VerifiedEmail = userInfo.VerifiedEmail
		if err := h.Storage.UpdateUser(&user); err != nil {
			log.Println("Failed to update user:", err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		log.Println(err.Error())
		return
	}

	frontendURL := utils.GetEnvVariable("FRONTEND_URL")
	redirectURL := fmt.Sprintf("%s/authentication/callback?token=%s", frontendURL, token)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
