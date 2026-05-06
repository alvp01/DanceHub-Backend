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

func TestHandlerLogout_SuccessAndRefreshRejected(t *testing.T) {
	_, service, jwtManager := setupAcademyTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	registerBody := map[string]any{
		"name":          "Logout Studio",
		"email":         "logout@example.com",
		"primary_phone": "955111222",
		"password":      "ABcdef123!^x",
	}
	registerRec := performJSONRequest(t, r, http.MethodPost, "/api/v1/academies/register", registerBody, "")
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 on register, got %d body=%s", registerRec.Code, registerRec.Body.String())
	}

	loginBody := map[string]any{
		"email":    "logout@example.com",
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

	logoutBody := map[string]any{"refresh_token": loginResp.RefreshToken}
	logoutRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/logout",
		logoutBody,
		"Bearer "+loginResp.AccessToken,
	)
	if logoutRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on logout, got %d body=%s", logoutRec.Code, logoutRec.Body.String())
	}

	refreshRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/refresh",
		map[string]any{"refresh_token": loginResp.RefreshToken},
		"",
	)
	if refreshRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on refresh after logout, got %d body=%s", refreshRec.Code, refreshRec.Body.String())
	}
}

func TestHandlerLogout_InvalidBodyAndNotFoundToken(t *testing.T) {
	_, service, jwtManager := setupAcademyTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	registerBody := map[string]any{
		"name":          "Token Studio",
		"email":         "token@example.com",
		"primary_phone": "955111223",
		"password":      "ABcdef123!^x",
	}
	registerRec := performJSONRequest(t, r, http.MethodPost, "/api/v1/academies/register", registerBody, "")
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 on register, got %d body=%s", registerRec.Code, registerRec.Body.String())
	}

	loginRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/login",
		map[string]any{"email": "token@example.com", "password": "ABcdef123!^x"},
		"",
	)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on login, got %d body=%s", loginRec.Code, loginRec.Body.String())
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decoding login response: %v", err)
	}

	invalidBodyRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/logout",
		map[string]any{},
		"Bearer "+loginResp.AccessToken,
	)
	if invalidBodyRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when refresh_token is missing, got %d body=%s", invalidBodyRec.Code, invalidBodyRec.Body.String())
	}

	firstLogoutRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/logout",
		map[string]any{"refresh_token": loginResp.RefreshToken},
		"Bearer "+loginResp.AccessToken,
	)
	if firstLogoutRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on first logout, got %d body=%s", firstLogoutRec.Code, firstLogoutRec.Body.String())
	}

	secondLogoutRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/logout",
		map[string]any{"refresh_token": loginResp.RefreshToken},
		"Bearer "+loginResp.AccessToken,
	)
	if secondLogoutRec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when token is already revoked, got %d body=%s", secondLogoutRec.Code, secondLogoutRec.Body.String())
	}
}

func TestHandlerLogout_RejectsAnotherAcademyToken(t *testing.T) {
	_, service, jwtManager := setupAcademyTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	firstReg := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/register",
		map[string]any{"name": "First Studio", "email": "first-handler@example.com", "primary_phone": "955111224", "password": "ABcdef123!^x"},
		"",
	)
	if firstReg.Code != http.StatusCreated {
		t.Fatalf("expected 201 for first academy register, got %d body=%s", firstReg.Code, firstReg.Body.String())
	}

	secondReg := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/register",
		map[string]any{"name": "Second Studio", "email": "second-handler@example.com", "primary_phone": "955111225", "password": "ABcdef123!^x"},
		"",
	)
	if secondReg.Code != http.StatusCreated {
		t.Fatalf("expected 201 for second academy register, got %d body=%s", secondReg.Code, secondReg.Body.String())
	}

	firstLogin := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/login",
		map[string]any{"email": "first-handler@example.com", "password": "ABcdef123!^x"},
		"",
	)
	if firstLogin.Code != http.StatusOK {
		t.Fatalf("expected 200 for first academy login, got %d body=%s", firstLogin.Code, firstLogin.Body.String())
	}

	secondLogin := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/login",
		map[string]any{"email": "second-handler@example.com", "password": "ABcdef123!^x"},
		"",
	)
	if secondLogin.Code != http.StatusOK {
		t.Fatalf("expected 200 for second academy login, got %d body=%s", secondLogin.Code, secondLogin.Body.String())
	}

	var firstLoginResp LoginResponse
	if err := json.Unmarshal(firstLogin.Body.Bytes(), &firstLoginResp); err != nil {
		t.Fatalf("decoding first login response: %v", err)
	}

	var secondLoginResp LoginResponse
	if err := json.Unmarshal(secondLogin.Body.Bytes(), &secondLoginResp); err != nil {
		t.Fatalf("decoding second login response: %v", err)
	}

	logoutRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/logout",
		map[string]any{"refresh_token": secondLoginResp.RefreshToken},
		"Bearer "+firstLoginResp.AccessToken,
	)
	if logoutRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when academy tries to logout another academy token, got %d body=%s", logoutRec.Code, logoutRec.Body.String())
	}

	refreshRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/refresh",
		map[string]any{"refresh_token": secondLoginResp.RefreshToken},
		"",
	)
	if refreshRec.Code != http.StatusOK {
		t.Fatalf("expected second academy refresh token to remain valid, got %d body=%s", refreshRec.Code, refreshRec.Body.String())
	}
}

func TestHandlerLogoutAll_RevokesAllSessions(t *testing.T) {
	_, service, jwtManager := setupAcademyTestDeps(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	h.RegisterRoutes(r, middleware.Auth(jwtManager))

	registerRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/register",
		map[string]any{"name": "All Sessions Studio", "email": "all-sessions@example.com", "primary_phone": "955111226", "password": "ABcdef123!^x"},
		"",
	)
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 on register, got %d body=%s", registerRec.Code, registerRec.Body.String())
	}

	loginReq := map[string]any{"email": "all-sessions@example.com", "password": "ABcdef123!^x"}
	firstLoginRec := performJSONRequest(t, r, http.MethodPost, "/api/v1/academies/login", loginReq, "")
	if firstLoginRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on first login, got %d body=%s", firstLoginRec.Code, firstLoginRec.Body.String())
	}

	secondLoginRec := performJSONRequest(t, r, http.MethodPost, "/api/v1/academies/login", loginReq, "")
	if secondLoginRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on second login, got %d body=%s", secondLoginRec.Code, secondLoginRec.Body.String())
	}

	var firstLoginResp LoginResponse
	if err := json.Unmarshal(firstLoginRec.Body.Bytes(), &firstLoginResp); err != nil {
		t.Fatalf("decoding first login response: %v", err)
	}

	var secondLoginResp LoginResponse
	if err := json.Unmarshal(secondLoginRec.Body.Bytes(), &secondLoginResp); err != nil {
		t.Fatalf("decoding second login response: %v", err)
	}

	logoutAllRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/logout-all",
		map[string]any{},
		"Bearer "+firstLoginResp.AccessToken,
	)
	if logoutAllRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on logout-all, got %d body=%s", logoutAllRec.Code, logoutAllRec.Body.String())
	}

	firstRefreshRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/refresh",
		map[string]any{"refresh_token": firstLoginResp.RefreshToken},
		"",
	)
	if firstRefreshRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for first revoked token, got %d body=%s", firstRefreshRec.Code, firstRefreshRec.Body.String())
	}

	secondRefreshRec := performJSONRequest(
		t,
		r,
		http.MethodPost,
		"/api/v1/academies/refresh",
		map[string]any{"refresh_token": secondLoginResp.RefreshToken},
		"",
	)
	if secondRefreshRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for second revoked token, got %d body=%s", secondRefreshRec.Code, secondRefreshRec.Body.String())
	}
}
