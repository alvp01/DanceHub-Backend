// internal/academy/handler.go
package academy

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/alvp01/DanceHub-Backend/internal/middleware"
	"github.com/alvp01/DanceHub-Backend/internal/validator"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	academies := router.Group("/api/v1/academies")
	academies.POST("/register", h.Register)
	academies.POST("/login", h.Login)
	academies.POST("/refresh", h.Refresh)

	protected := academies.Group("")
	protected.Use(authMiddleware)
	protected.POST("/logout", h.Logout)
	protected.POST("/logout-all", h.LogoutAll)
	protected.GET("/me", h.Me)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "cuerpo de la petición inválido")
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeError(c, 401, err.Error())
			return
		}
		writeError(c, 500, "error interno del servidor")
		return
	}

	writeJSON(c, 200, resp)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "cuerpo de la petición inválido")
		return
	}

	resp, err := h.service.RefreshTokens(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeError(c, 401, "token inválido o expirado")
			return
		}
		writeError(c, 500, "error interno del servidor")
		return
	}

	writeJSON(c, 200, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "cuerpo de la petición inválido")
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			writeError(c, 404, "token no encontrado")
			return
		}
		writeError(c, 500, "error interno del servidor")
		return
	}

	writeJSON(c, 200, map[string]string{"message": "sesión cerrada correctamente"})
}

func (h *Handler) LogoutAll(c *gin.Context) {
	// Obtener academy_id desde el JWT (ya validado por el middleware)
	claims, ok := middleware.GetClaims(c)
	if !ok {
		writeError(c, 401, "no autorizado")
		return
	}

	if err := h.service.LogoutAll(c.Request.Context(), claims.AcademyID); err != nil {
		writeError(c, 500, "error interno del servidor")
		return
	}

	writeJSON(c, 200, map[string]string{
		"message": "todas las sesiones cerradas correctamente",
	})
}

func (h *Handler) Me(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		writeError(c, 401, "no autorizado")
		return
	}

	writeJSON(c, 200, map[string]any{
		"academy_id":   claims.AcademyID,
		"academy_name": claims.AcademyName,
		"token_type":   claims.TokenType,
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "cuerpo de la petición inválido")
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		var pwErr *validator.PasswordValidationError
		if errors.As(err, &pwErr) {
			writeJSON(c, 422, map[string]any{
				"error":  "password inválido",
				"detail": pwErr.Errors,
			})
			return
		}

		switch {
		case errors.Is(err, ErrEmailAlreadyExists),
			errors.Is(err, ErrNameAlreadyExists),
			errors.Is(err, ErrPhoneAlreadyExists):
			writeError(c, 409, err.Error())
			return
		}

		writeError(c, 500, "error interno del servidor")
		return
	}

	writeJSON(c, 201, resp)
}

func writeJSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func writeError(c *gin.Context, status int, message string) {
	writeJSON(c, status, map[string]string{"error": message})
}
