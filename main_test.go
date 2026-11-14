package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, path string) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec.Code, rec.Body.String()
}

func TestRootReportsSetAndUnsetWithoutValues(t *testing.T) {
	t.Setenv("ORDERS_API_KEY", "super-secret-value")
	t.Setenv("ORDERS_GREETING", "")
	code, body := get(t, "/")
	if code != http.StatusOK {
		t.Fatalf("status %d", code)
	}
	if !strings.Contains(body, `"ORDERS_API_KEY":"set"`) || !strings.Contains(body, `"ORDERS_GREETING":"unset"`) {
		t.Fatalf("unexpected body %s", body)
	}
	if strings.Contains(body, "super-secret-value") {
		t.Fatal("the response must never contain a secret value")
	}
}

func TestHealthz(t *testing.T) {
	if code, body := get(t, "/healthz"); code != http.StatusOK || body != "ok" {
		t.Fatalf("got %d %q", code, body)
	}
}

func TestUnknownPathIs404(t *testing.T) {
	if code, _ := get(t, "/nope"); code != http.StatusNotFound {
		t.Fatalf("got %d", code)
	}
}
