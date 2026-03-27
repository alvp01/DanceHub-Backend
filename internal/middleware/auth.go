// internal/middleware/auth.go
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	jwtpkg "github.com/alvp01/DanceHub-Backend/internal/jwt"
)

// ContextKey tipo propio para evitar colisiones en el contexto
type contextKey string

const AcademyClaimsKey contextKey = "academy_claims"

func Auth(jwtManager *jwtpkg.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extraer el token del header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "authorization header requerido")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeError(w, http.StatusUnauthorized, "formato inválido: Bearer <token>")
				return
			}

			tokenStr := parts[1]

			// 2. Validar el access token
			claims, err := jwtManager.ValidateAccessToken(tokenStr)
			if err != nil {
				if errors.Is(err, jwtpkg.ErrExpiredToken) {
					writeError(w, http.StatusUnauthorized, "token expirado")
					return
				}
				writeError(w, http.StatusUnauthorized, "token inválido")
				return
			}

			// 3. Inyectar claims en el contexto
			ctx := context.WithValue(r.Context(), AcademyClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims extrae los claims del contexto de forma segura
func GetClaims(r *http.Request) (*jwtpkg.Claims, bool) {
	claims, ok := r.Context().Value(AcademyClaimsKey).(*jwtpkg.Claims)
	return claims, ok
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + msg + `"}`))
}
