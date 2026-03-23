package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGatewayMiddleware(t *testing.T) {
	mws := GatewayMiddleware()
	if len(mws) == 0 {
		t.Error("expected at least one middleware")
	}

	found := false
	for _, mw := range mws {
		// We can't easily compare function pointers in Go,
		// but we can check if it works as expected.
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		panicHandler := func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			panic("middleware panic")
		}

		handler := mw(panicHandler)
		handler(w, req, nil)

		if w.Code == http.StatusBadGateway {
			found = true
			break
		}
	}

	if !found {
		t.Error("RecoverMiddleware not found in GatewayMiddleware")
	}
}
