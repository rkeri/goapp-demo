package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func resetConfig(t *testing.T) {
	t.Helper()
	configFile = filepath.Join(t.TempDir(), "config.json")
	config = map[string]string{}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestVersionHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	versionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), version) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestEnvHandler(t *testing.T) {
	t.Setenv("ENVIRONMENT", "test")

	req := httptest.NewRequest(http.MethodGet, "/env", nil)
	rec := httptest.NewRecorder()

	envHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"environment":"test"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestConfigCreateAndGet(t *testing.T) {
	resetConfig(t)

	body := strings.NewReader(`{"name":"database_url","value":"postgres://example"}`)
	req := httptest.NewRequest(http.MethodPost, "/config", body)
	rec := httptest.NewRecorder()

	createConfigHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var created ConfigItem
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if created.Name != "database_url" || created.Value != "postgres://example" {
		t.Fatalf("unexpected response: %+v", created)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/config/database_url", nil)
	getRec := httptest.NewRecorder()

	configItemHandler(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getRec.Code)
	}
	if !strings.Contains(getRec.Body.String(), "postgres://example") {
		t.Fatalf("unexpected body: %s", getRec.Body.String())
	}
}

func TestConfigGetMissing(t *testing.T) {
	resetConfig(t)

	req := httptest.NewRequest(http.MethodGet, "/config/does-not-exist", nil)
	rec := httptest.NewRecorder()

	configItemHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestConfigDelete(t *testing.T) {
	resetConfig(t)
	config["database_url"] = "postgres://example"

	delReq := httptest.NewRequest(http.MethodDelete, "/config/database_url", nil)
	delRec := httptest.NewRecorder()

	configItemHandler(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", delRec.Code)
	}
	if !strings.Contains(delRec.Body.String(), `"deleted":true`) {
		t.Fatalf("unexpected body: %s", delRec.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/config/database_url", nil)
	getRec := httptest.NewRecorder()

	configItemHandler(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", getRec.Code)
	}
}
