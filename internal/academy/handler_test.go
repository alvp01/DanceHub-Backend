package academy

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alvp01/DanceHub-Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestHandlerRegisterLoginAndMe(t *testing.T) {
	_, service, jwtManager := setupAcademyTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	registerBody := map[string]any{
		"name":          "Studio One",
		"email":         "studio@example.com",
		"primary_phone": "988777666",
		"password":      "ABcdef123!^x",
	}
	registerRec := performJSONRequest(t, r, http.MethodPost, "/api/v1/academies/register", registerBody, "")
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 on register, got %d body=%s", registerRec.Code, registerRec.Body.String())
	}

	loginBody := map[string]any{
		"email":    "studio@example.com",
		"password": "ABcdef123!^x",
	}
	loginRec := performJSONRequest(t, r, http.MethodPost, "/api/v1/academies/login", loginBody, "")
	if loginRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on login, got %d body=%s", loginRec.Code, loginRec.Body.String())
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decoding login response: %v", err)
	}

	if loginResp.AccessToken == "" {
		t.Fatalf("expected access token in login response")
	}

	meRec := performJSONRequest(
		t,
		r,
		http.MethodGet,
		"/api/v1/academies/me",
		nil,
		"Bearer "+loginResp.AccessToken,
	)
	if meRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on me endpoint, got %d body=%s", meRec.Code, meRec.Body.String())
	}

	var meResp map[string]any
	if err := json.Unmarshal(meRec.Body.Bytes(), &meResp); err != nil {
		t.Fatalf("decoding me response: %v", err)
	}
	if meResp["academy_id"] == "" {
		t.Fatalf("expected academy_id in me response")
	}
}

func TestHandlerProtectedRouteWithoutAuth(t *testing.T) {
	_, service, jwtManager := setupAcademyTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	logoutAllRec := performJSONRequest(t, r, http.MethodPost, "/api/v1/academies/logout-all", map[string]any{}, "")
	if logoutAllRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on protected route without auth, got %d", logoutAllRec.Code)
	}
}
