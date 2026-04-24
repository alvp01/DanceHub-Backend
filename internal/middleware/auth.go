// internal/middleware/auth.go
package middleware

import (
	"errors"
	"strings"

	jwtpkg "github.com/alvp01/DanceHub-Backend/internal/jwt"
	"github.com/gin-gonic/gin"
)

const AcademyClaimsKey = "academy_claims"

func Auth(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			writeError(c, 401, "authorization header requerido")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(c, 401, "formato inválido: Bearer <token>")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims, err := jwtManager.ValidateAccessToken(tokenStr)
		if err != nil {
			if errors.Is(err, jwtpkg.ErrExpiredToken) {
				writeError(c, 401, "token expirado")
				c.Abort()
				return
			}
			writeError(c, 401, "token inválido")
			c.Abort()
			return
		}

		c.Set(AcademyClaimsKey, claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) (*jwtpkg.Claims, bool) {
	value, ok := c.Get(AcademyClaimsKey)
	if !ok {
		return nil, false
	}

	claims, ok := value.(*jwtpkg.Claims)
	return claims, ok
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
