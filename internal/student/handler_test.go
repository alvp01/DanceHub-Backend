package student

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alvp01/DanceHub-Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func setupStudentHandlerTestRouter(t *testing.T) (*Repository, *gin.Engine, string) {
	t.Helper()

	repo, service, jwtManager := setupStudentTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	accessToken, err := jwtManager.GenerateAccessToken("academy-1", "Academy One")
	if err != nil {
		t.Fatalf("generating access token: %v", err)
	}

	return repo, r, "Bearer " + accessToken
}

func studentCreateBody(email, phone, idDocument string) map[string]any {
	return map[string]any{
		"academy_id":  "academy-1",
		"name":        "Laura",
		"last_name":   "Gomez",
		"email":       email,
		"phone":       phone,
		"id_document": idDocument,
		"birth_date":  "2001-03-03",
		"address":     "Street 30",
		"allergies":   "none",
		"pathologies": "none",
	}
}

func TestHandlerStudentRoutesRequireAuth(t *testing.T) {
	_, service, jwtManager := setupStudentTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	tests := []struct {
		name   string
		method string
		path   string
		body   any
	}{
		{name: "create", method: http.MethodPost, path: "/api/v1/academy/students/register", body: map[string]any{}},
		{name: "find all", method: http.MethodGet, path: "/api/v1/academy/students/", body: nil},
		{name: "find by id document", method: http.MethodPost, path: "/api/v1/academy/students/", body: map[string]any{"id_document": "DOC-300"}},
		{name: "update", method: http.MethodPut, path: "/api/v1/academy/students/", body: map[string]any{}},
		{name: "delete", method: http.MethodDelete, path: "/api/v1/academy/students/any-id", body: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performJSONRequest(t, r, tt.method, tt.path, tt.body, "")
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401 without auth, got %d body=%s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandlerCreateStudentWithAuth(t *testing.T) {
	_, r, authHeader := setupStudentHandlerTestRouter(t)

	createRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academy/students/register",
		studentCreateBody("laura@example.com", "300555111", "DOC-400"),
		authHeader,
	)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 on create, got %d body=%s", createRec.Code, createRec.Body.String())
	}
}

func TestHandlerFindByDocumentWithAuth(t *testing.T) {
	repo, r, authHeader := setupStudentHandlerTestRouter(t)
	ctx := context.Background()

	if err := repo.Create(ctx, &Student{
		AcademyID:   "academy-1",
		Name:        "Laura",
		LastName:    "Gomez",
		Email:       "laura.find@example.com",
		Phone:       "300555121",
		IdDocument:  "DOC-410",
		BirthDate:   "2001-03-03",
		Address:     "Street 30",
		Allergies:   "none",
		Pathologies: "none",
	}); err != nil {
		t.Fatalf("creating fixture student: %v", err)
	}

	findByDocRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academy/students/",
		map[string]any{"id_document": "DOC-410"},
		authHeader,
	)
	if findByDocRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on find by document, got %d body=%s", findByDocRec.Code, findByDocRec.Body.String())
	}

	var found Student
	if err := json.Unmarshal(findByDocRec.Body.Bytes(), &found); err != nil {
		t.Fatalf("decoding find by document response: %v", err)
	}
	if found.ID == "" {
		t.Fatalf("expected found student id")
	}
}

func TestHandlerFindAllWithAuth(t *testing.T) {
	repo, r, authHeader := setupStudentHandlerTestRouter(t)
	ctx := context.Background()

	if err := repo.Create(ctx, &Student{
		AcademyID:   "academy-1",
		Name:        "Laura",
		LastName:    "Gomez",
		Email:       "laura.list@example.com",
		Phone:       "300555122",
		IdDocument:  "DOC-411",
		BirthDate:   "2001-03-03",
		Address:     "Street 31",
		Allergies:   "none",
		Pathologies: "none",
	}); err != nil {
		t.Fatalf("creating fixture student: %v", err)
	}

	findAllRec := performJSONRequest(t, r, http.MethodGet, "/api/v1/academy/students/", nil, authHeader)
	if findAllRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on find all, got %d body=%s", findAllRec.Code, findAllRec.Body.String())
	}

	var all []Student
	if err := json.Unmarshal(findAllRec.Body.Bytes(), &all); err != nil {
		t.Fatalf("decoding find all response: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 student in response, got %d", len(all))
	}
}

func TestHandlerUpdateStudentWithAuth(t *testing.T) {
	repo, r, authHeader := setupStudentHandlerTestRouter(t)
	ctx := context.Background()

	fixture := &Student{
		AcademyID:   "academy-1",
		Name:        "Laura",
		LastName:    "Gomez",
		Email:       "laura.update@example.com",
		Phone:       "300555123",
		IdDocument:  "DOC-412",
		BirthDate:   "2001-03-03",
		Address:     "Street 32",
		Allergies:   "none",
		Pathologies: "none",
	}
	if err := repo.Create(ctx, fixture); err != nil {
		t.Fatalf("creating fixture student: %v", err)
	}

	updateBody := map[string]any{
		"id":          fixture.ID,
		"name":        "Laura Maria",
		"last_name":   "Gomez",
		"email":       "laura.maria@example.com",
		"phone":       "300555124",
		"id_document": "DOC-413",
		"birth_date":  "2001-03-03",
		"address":     "Street 33",
		"allergies":   "pollen",
		"pathologies": "none",
	}
	updateRec := performJSONRequest(t, r, http.MethodPut, "/api/v1/academy/students/", updateBody, authHeader)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d body=%s", updateRec.Code, updateRec.Body.String())
	}
}

func TestHandlerDeleteStudentWithAuth(t *testing.T) {
	repo, r, authHeader := setupStudentHandlerTestRouter(t)
	ctx := context.Background()

	fixture := &Student{
		AcademyID:   "academy-1",
		Name:        "Laura",
		LastName:    "Gomez",
		Email:       "laura.delete@example.com",
		Phone:       "300555125",
		IdDocument:  "DOC-414",
		BirthDate:   "2001-03-03",
		Address:     "Street 34",
		Allergies:   "none",
		Pathologies: "none",
	}
	if err := repo.Create(ctx, fixture); err != nil {
		t.Fatalf("creating fixture student: %v", err)
	}

	deleteRec := performJSONRequest(t, r, http.MethodDelete, "/api/v1/academy/students/"+fixture.ID, nil, authHeader)
	if deleteRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d body=%s", deleteRec.Code, deleteRec.Body.String())
	}
}
