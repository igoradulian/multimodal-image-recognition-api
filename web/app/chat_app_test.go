package app

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetBotResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/responses", GetBotResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/responses", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload["id"] == "" || payload["text"] == "" {
		t.Fatalf("expected id and text in response, got %#v", payload)
	}
}

func TestPostOpenAIVisionMissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/openai/vision", PostOpenAIVision)

	body, contentType := newMultipartBody(t, func(_ *multipart.Writer) {})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/openai/vision", body)
	req.Header.Set("Content-Type", contentType)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGoogleChatServiceMissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/google/vision", GoogleChatService)

	body, contentType := newMultipartBody(t, func(_ *multipart.Writer) {})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/google/vision", body)
	req.Header.Set("Content-Type", contentType)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPostQuestionMissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/responses", PostQuestion)

	body, contentType := newMultipartBody(t, func(w *multipart.Writer) {
		if err := w.WriteField("message", "hello"); err != nil {
			t.Fatalf("failed to write form field: %v", err)
		}
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/responses", body)
	req.Header.Set("Content-Type", contentType)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func newMultipartBody(t *testing.T, write func(*multipart.Writer)) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	write(writer)
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}
	return body, writer.FormDataContentType()
}
