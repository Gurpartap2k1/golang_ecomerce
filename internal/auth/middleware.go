package auth

import (
	"net/http"
	"strings"
)

func (m *JwtManager) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		const bearer = "Bearer "

		if !strings.HasPrefix(authHeader, bearer) {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearer)

		claims, err := m.Verify(tokenString)

		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := WithUserId(r.Context(), claims.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
