package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shadow-nova/backend/internal/database"
)

// mockDB implements the minimal database.Service interface needed for testing
type mockIdempotencyDB struct {
	database.MockService
	cache map[string]*database.IdempotentResponse
}

func newMockIdempotencyDB() *mockIdempotencyDB {
	return &mockIdempotencyDB{
		cache: make(map[string]*database.IdempotentResponse),
	}
}

func (m *mockIdempotencyDB) StoreIdempotentResponse(ctx context.Context, key string, userID int, path, method string, status int, body string, expiresAt time.Time) error {
	m.cache[key] = &database.IdempotentResponse{
		Key:        key,
		StatusCode: status,
		Body:       body,
	}
	return nil
}

func (m *mockIdempotencyDB) GetIdempotentResponse(ctx context.Context, key string, userID int) (*database.IdempotentResponse, error) {
	if resp, ok := m.cache[key]; ok {
		return resp, nil
	}
	return nil, nil
}

func TestIdempotencyMiddleware_NoKey(t *testing.T) {
	db := newMockIdempotencyDB()
	handler := Idempotency(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))

	req := httptest.NewRequest("POST", "/api/test", nil)
	req = req.WithContext(context.WithValue(req.Context(), UserIDKey, 1))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Should not be cached
	if len(db.cache) != 0 {
		t.Error("Request without idempotency key should not be cached")
	}
}

func TestIdempotencyMiddleware_WithKey(t *testing.T) {
	db := newMockIdempotencyDB()
	callCount := 0
	handler := Idempotency(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))

	idempotencyKey := "test-key-123"

	// First request
	req1 := httptest.NewRequest("POST", "/api/test", nil)
	req1.Header.Set("Idempotency-Key", idempotencyKey)
	req1 = req1.WithContext(context.WithValue(req1.Context(), UserIDKey, 1))
	w1 := httptest.NewRecorder()

	handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}
	if callCount != 1 {
		t.Errorf("Expected handler to be called once, got %d", callCount)
	}

	// Second request with same key
	req2 := httptest.NewRequest("POST", "/api/test", nil)
	req2.Header.Set("Idempotency-Key", idempotencyKey)
	req2 = req2.WithContext(context.WithValue(req2.Context(), UserIDKey, 1))
	w2 := httptest.NewRecorder()

	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
	if callCount != 1 {
		t.Errorf("Expected handler to be called only once, got %d", callCount)
	}

	// Check for replay header
	if w2.Header().Get("X-Idempotent-Replay") != "true" {
		t.Error("Expected X-Idempotent-Replay header on cached response")
	}

	// Responses should be identical
	if w1.Body.String() != w2.Body.String() {
		t.Error("Cached response body should match original")
	}
}

func TestIdempotencyMiddleware_OnlyMutatingMethods(t *testing.T) {
	db := newMockIdempotencyDB()
	handler := Idempotency(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	methods := []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}
	shouldCache := map[string]bool{
		"POST":  true,
		"PUT":   true,
		"PATCH": true,
	}

	for _, method := range methods {
		db.cache = make(map[string]*database.IdempotentResponse) // Clear cache
		req := httptest.NewRequest(method, "/api/test", nil)
		req.Header.Set("Idempotency-Key", "test-key")
		req = req.WithContext(context.WithValue(req.Context(), UserIDKey, 1))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		cached := len(db.cache) > 0
		if cached != shouldCache[method] {
			t.Errorf("Method %s: expected cached=%v, got %v", method, shouldCache[method], cached)
		}
	}
}

func TestIdempotencyMiddleware_DifferentUsers(t *testing.T) {
	db := newMockIdempotencyDB()
	callCount := 0
	handler := Idempotency(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))

	idempotencyKey := "shared-key"

	// User 1
	req1 := httptest.NewRequest("POST", "/api/test", nil)
	req1.Header.Set("Idempotency-Key", idempotencyKey)
	req1 = req1.WithContext(context.WithValue(req1.Context(), UserIDKey, 1))
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	// User 2 with same key
	req2 := httptest.NewRequest("POST", "/api/test", nil)
	req2.Header.Set("Idempotency-Key", idempotencyKey)
	req2 = req2.WithContext(context.WithValue(req2.Context(), UserIDKey, 2))
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	// Both should execute (different users)
	if callCount != 2 {
		t.Errorf("Expected handler to be called twice (once per user), got %d", callCount)
	}
}

func TestIdempotencyMiddleware_NoUserContext(t *testing.T) {
	db := newMockIdempotencyDB()
	callCount := 0
	handler := Idempotency(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set("Idempotency-Key", "test-key")
	// No user context
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if callCount != 1 {
		t.Errorf("Expected handler to be called, got %d", callCount)
	}
	// Should not be cached without user context
	if len(db.cache) != 0 {
		t.Error("Request without user context should not be cached")
	}
}

func TestIdempotencyMiddleware_CapturesStatusCode(t *testing.T) {
	db := newMockIdempotencyDB()
	handler := Idempotency(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":123}`))
	}))

	idempotencyKey := "test-key-201"

	// First request
	req1 := httptest.NewRequest("POST", "/api/test", nil)
	req1.Header.Set("Idempotency-Key", idempotencyKey)
	req1 = req1.WithContext(context.WithValue(req1.Context(), UserIDKey, 1))
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w1.Code)
	}

	// Second request should return same status code
	req2 := httptest.NewRequest("POST", "/api/test", nil)
	req2.Header.Set("Idempotency-Key", idempotencyKey)
	req2 = req2.WithContext(context.WithValue(req2.Context(), UserIDKey, 1))
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusCreated {
		t.Errorf("Expected cached status 201, got %d", w2.Code)
	}
	if w2.Body.String() != `{"id":123}` {
		t.Errorf("Expected cached body, got %s", w2.Body.String())
	}
}

func TestIdempotencyMiddleware_LargeResponse(t *testing.T) {
	db := newMockIdempotencyDB()
	largeBody := bytes.Repeat([]byte("x"), 10000)
	handler := Idempotency(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(largeBody)
	}))

	idempotencyKey := "test-large"

	req1 := httptest.NewRequest("POST", "/api/test", nil)
	req1.Header.Set("Idempotency-Key", idempotencyKey)
	req1 = req1.WithContext(context.WithValue(req1.Context(), UserIDKey, 1))
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	req2 := httptest.NewRequest("POST", "/api/test", nil)
	req2.Header.Set("Idempotency-Key", idempotencyKey)
	req2 = req2.WithContext(context.WithValue(req2.Context(), UserIDKey, 1))
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if !bytes.Equal(w1.Body.Bytes(), w2.Body.Bytes()) {
		t.Error("Large cached response should match original")
	}
	if len(w2.Body.Bytes()) != 10000 {
		t.Errorf("Expected 10000 bytes, got %d", len(w2.Body.Bytes()))
	}
}
