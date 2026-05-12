package student

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	students := router.Group("/api/v1/academy/students")
	students.Use(authMiddleware)
	students.POST("/register", h.Create)
	students.GET("/", h.FindAll)
	students.POST("/", h.FindByIdDocument)
	students.PUT("/", h.Update)
	students.DELETE("/:id", h.Delete)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "cuerpo de la petición inválido")
		return
	}

	student, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists),
			errors.Is(err, ErrPhoneAlreadyExists),
			errors.Is(err, ErrIdDocAlreadyExists):
			writeError(c, 400, err.Error())
		default:
			writeError(c, 500, "error interno del servidor")
		}
		return
	}

	writeJSON(c, 201, student)
}

func (h *Handler) FindAll(c *gin.Context) {
	students, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		writeError(c, 500, "error interno del servidor")
		return
	}

	writeJSON(c, 200, students)
}

func (h *Handler) FindByIdDocument(c *gin.Context) {
	var req FindByIdDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "cuerpo de la petición inválido")
		return
	}

	student, err := h.service.FindByIdDocument(c.Request.Context(), req.IdDocument)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			writeError(c, 404, err.Error())
		default:
			writeError(c, 500, "error interno del servidor")
		}
		return
	}

	writeJSON(c, 200, student)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "cuerpo de la petición inválido")
		return
	}

	student, err := h.service.Update(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			writeError(c, 404, err.Error())
		case errors.Is(err, ErrEmailAlreadyExists),
			errors.Is(err, ErrPhoneAlreadyExists),
			errors.Is(err, ErrIdDocAlreadyExists):
			writeError(c, 400, err.Error())
		default:
			writeError(c, 500, "error interno del servidor")
		}
		return
	}

	writeJSON(c, 200, student)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		writeError(c, 400, "el ID del estudiante es obligatorio")
		return
	}

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			writeError(c, 404, err.Error())
		default:
			writeError(c, 500, "error interno del servidor")
		}
		return
	}

	writeJSON(c, 200, map[string]string{"message": "estudiante eliminado exitosamente"})
}

func writeError(c *gin.Context, status int, message string) {
	writeJSON(c, status, map[string]string{"error": message})
}

func writeJSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}
