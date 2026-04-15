// internal/academy/handler.go
package academy

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alvp01/DanceHub-Backend/internal/middleware"
	"github.com/alvp01/DanceHub-Backend/internal/validator"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.HandleFunc("POST /api/v1/academies/register", h.Register)
	mux.HandleFunc("POST /api/v1/academies/login", h.Login)
	mux.HandleFunc("POST /api/v1/academies/refresh", h.Refresh)

	mux.Handle("POST /api/v1/academies/logout",
		authMiddleware(http.HandlerFunc(h.Logout)))

	mux.Handle("POST /api/v1/academies/logout-all",
		authMiddleware(http.HandlerFunc(h.LogoutAll)))

	mux.Handle("GET /api/v1/academies/me",
		authMiddleware(http.HandlerFunc(h.Me)))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
		return
	}
	defer r.Body.Close()

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "error interno del servidor")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
		return
	}
	defer r.Body.Close()

	resp, err := h.service.RefreshTokens(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "token inválido o expirado")
			return
		}
		writeError(w, http.StatusInternalServerError, "error interno del servidor")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
		return
	}
	defer r.Body.Close()

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			writeError(w, http.StatusNotFound, "token no encontrado")
			return
		}
		writeError(w, http.StatusInternalServerError, "error interno del servidor")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "sesión cerrada correctamente"})
}

func (h *Handler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	// Obtener academy_id desde el JWT (ya validado por el middleware)
	claims, ok := middleware.GetClaims(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "no autorizado")
		return
	}

	if err := h.service.LogoutAll(r.Context(), claims.AcademyID); err != nil {
		writeError(w, http.StatusInternalServerError, "error interno del servidor")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "todas las sesiones cerradas correctamente",
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "no autorizado")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"academy_id":   claims.AcademyID,
		"academy_name": claims.AcademyName,
		"token_type":   claims.TokenType,
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
		return
	}
	defer r.Body.Close()

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		var pwErr *validator.PasswordValidationError
		if errors.As(err, &pwErr) {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
				"error":  "password inválido",
				"detail": pwErr.Errors,
			})
			return
		}

		switch {
		case errors.Is(err, ErrEmailAlreadyExists),
			errors.Is(err, ErrNameAlreadyExists),
			errors.Is(err, ErrPhoneAlreadyExists):
			writeError(w, http.StatusConflict, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, "error interno del servidor")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
