package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc( func (w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }
		
		// Check if header has Bearer type of auth
		parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }
		
		// Take only second part of header(token value)
		tokenString := parts[1]

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}

		userID := strconv.FormatFloat(claims["userID"].(float64), 'f', -1, 64)
		r.Header.Set("userID", userID)

		next.ServeHTTP(w, r)
	})
}
