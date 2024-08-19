package api

import (
	"net/http"
	"strings"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			sendErrorResponse(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if header has Bearer type of auth
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendErrorResponse(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Take only second part of header(token value)
		tokenString := parts[1]

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			sendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID := claims["userID"].(string)
		r.Header.Set("userID", userID)

		next.ServeHTTP(w, r)
	})
}
