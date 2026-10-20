package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPGenerationProviderUsesConfiguredEndpointAndModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("request method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization header = %q, want bearer test-key", got)
		}
		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		if body.Model != "configured-model" {
			t.Errorf("model = %q, want configured-model", body.Model)
		}
		if len(body.Messages) != 1 || body.Messages[0].Role != "user" || body.Messages[0].Content != "test prompt" {
			t.Errorf("messages = %+v, want one user message with the prompt", body.Messages)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"generated"}}]}`))
	}))
	defer server.Close()

	provider, err := NewHTTPGenerationProvider(HTTPGenerationProviderConfig{
		Endpoint: server.URL,
		APIKey:   "test-key",
		Model:    "configured-model",
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := provider.Generate(context.Background(), "test prompt")
	if err != nil {
		t.Fatal(err)
	}
	if got != "generated" {
		t.Errorf("Generate() = %q, want generated", got)
	}
}

func TestHTTPGenerationProviderAllowsNoModelOrAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization header = %q, want none", got)
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		if _, exists := body["model"]; exists {
			t.Errorf("request unexpectedly includes model: %s", body["model"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"generated"}}]}`))
	}))
	defer server.Close()

	provider, err := NewHTTPGenerationProvider(HTTPGenerationProviderConfig{Endpoint: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := provider.Generate(context.Background(), "test prompt"); err != nil {
		t.Fatal(err)
	}
}

func TestNewHTTPGenerationProviderRejectsInvalidEndpoint(t *testing.T) {
	if _, err := NewHTTPGenerationProvider(HTTPGenerationProviderConfig{Endpoint: "not-a-url"}); err == nil {
		t.Fatal("NewHTTPGenerationProvider succeeded with an invalid endpoint")
	}
}
