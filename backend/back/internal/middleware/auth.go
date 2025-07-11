package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"back/internal/utils"

	"github.com/gorilla/websocket"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Token Format", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invallid Token", http.StatusUnauthorized)
			return
		}
		log.Println(claims)

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "created_at", claims.CreatedAt)
		
		log.Println("Context: ", ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// middleware for endpoints initiating a ws connection
func WSAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		protocols := websocket.Subprotocols(r)
		var token string

		for _, protocol := range protocols {
			if strings.HasPrefix(protocol, "auth.") {
				token = strings.TrimPrefix(protocol, "auth.")
				break
			}
		}

		if token == "" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}


		claims, err := utils.ValidateJWT(token)
		if err != nil {
			log.Println("ERROR validating token: ", err)
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		log.Println("Claims: ", claims)

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "created_at", claims.CreatedAt)

		log.Println("Context: ", ctx)

		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
