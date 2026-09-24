package orders

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerCreatesAuthorizedValidOrder(t *testing.T) {
	handler := NewHandler("demo-key", NewMemoryStore())
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"customer":"Ada","item":"keyboard","quantity":2}`))
	request.Header.Set("Authorization", "Bearer demo-key")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusCreated, response.Body.String())
	}
}

func TestHandlerRejectsUnauthorizedRequest(t *testing.T) {
	handler := NewHandler("demo-key", NewMemoryStore())
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"customer":"Ada","item":"keyboard","quantity":2}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestHandlerRejectsInvalidOrder(t *testing.T) {
	handler := NewHandler("demo-key", NewMemoryStore())
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"customer":"Ada","item":"keyboard","quantity":0}`))
	request.Header.Set("Authorization", "Bearer demo-key")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
