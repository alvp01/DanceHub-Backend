package student

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtpkg "github.com/alvp01/DanceHub-Backend/internal/jwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupStudentTestDeps(t *testing.T) (*Repository, *Service, *jwtpkg.Manager) {
	t.Helper()

	t.Setenv("JWT_ACCESS_SECRET", "student-test-access-secret")
	t.Setenv("JWT_REFRESH_SECRET", "student-test-refresh-secret")
	t.Setenv("JWT_ACCESS_DURATION_MINUTES", "15")
	t.Setenv("JWT_REFRESH_DURATION_DAYS", "7")

	dsn := fmt.Sprintf("file:%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&Student{}); err != nil {
		t.Fatalf("migrating models: %v", err)
	}

	jwtManager, err := jwtpkg.NewManager()
	if err != nil {
		t.Fatalf("creating jwt manager: %v", err)
	}

	repo := NewRepository(db)
	service := NewService(repo)

	return repo, service, jwtManager
}

func performJSONRequest(
	t *testing.T,
	r http.Handler,
	method string,
	path string,
	body any,
	authorization string,
) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshaling request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}
