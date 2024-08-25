package handlers

import (
	"net/http"
	"strings"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
)

func (h *DefaultHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.Storage.GetUsers()
	if err != nil {
		utils.SendErrorResponse(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	utils.SendDataResponse(w, "Users retrieved successfully", http.StatusOK, users)

}
