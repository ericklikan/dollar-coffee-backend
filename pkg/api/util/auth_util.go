package util

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response map[string]interface{}
		tokenHeader := r.Header.Get("Authorization")

		// token header: `Bearer {token-body}`
		splitted := strings.Split(tokenHeader, " ")
		if tokenHeader == "" || len(splitted) != 2 {
			response = Message("Missing/Invalid/Malformed auth token")
			Respond(w, http.StatusForbidden, response)
			return
		}

		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(splitted[1], tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil || !token.Valid {
			response = Message("Something was wrong with auth token")
			Respond(w, http.StatusForbidden, response)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		ctx = context.WithValue(ctx, "role", tk.Role)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
