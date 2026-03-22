package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware(t *testing.T) {
	panicHandler := func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		panic("test panic")
	}

	middleware := RecoverMiddleware(panicHandler)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	middleware(w, req, nil)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status code %d, got %d", http.StatusBadGateway, resp.StatusCode)
	}

	expectedBody := http.StatusText(http.StatusBadGateway)
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
	}
}
