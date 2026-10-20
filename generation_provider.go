package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	generationProviderEndpointEnv = "LLM_API_URL"
	generationProviderAPIKeyEnv   = "LLM_API_KEY"
	generationProviderModelEnv    = "LLM_MODEL"
)

type HTTPGenerationProviderConfig struct {
	Endpoint string
	APIKey   string
	Model    string
	Client   *http.Client
}

type HTTPGenerationProvider struct {
	endpoint string
	apiKey   string
	model    string
	client   *http.Client
}

func HTTPGenerationProviderFromEnv() (*HTTPGenerationProvider, error) {
	endpoint := strings.TrimSpace(os.Getenv(generationProviderEndpointEnv))
	if endpoint == "" {
		return nil, fmt.Errorf("%s environment variable is required for llm steps", generationProviderEndpointEnv)
	}
	return NewHTTPGenerationProvider(HTTPGenerationProviderConfig{
		Endpoint: endpoint,
		APIKey:   os.Getenv(generationProviderAPIKeyEnv),
		Model:    os.Getenv(generationProviderModelEnv),
	})
}

func NewHTTPGenerationProvider(config HTTPGenerationProviderConfig) (*HTTPGenerationProvider, error) {
	endpoint := strings.TrimSpace(config.Endpoint)
	parsedEndpoint, err := url.ParseRequestURI(endpoint)
	if err != nil || parsedEndpoint.Scheme == "" || parsedEndpoint.Host == "" || (parsedEndpoint.Scheme != "http" && parsedEndpoint.Scheme != "https") {
		return nil, fmt.Errorf("invalid llm endpoint %q: expected an absolute HTTP or HTTPS URL", endpoint)
	}
	client := config.Client
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPGenerationProvider{
		endpoint: endpoint,
		apiKey:   strings.TrimSpace(config.APIKey),
		model:    strings.TrimSpace(config.Model),
		client:   client,
	}, nil
}

func (p *HTTPGenerationProvider) Generate(ctx context.Context, prompt string) (string, error) {
	requestBody, err := json.Marshal(struct {
		Model    string `json:"model,omitempty"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}{
		Model: p.model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return "", fmt.Errorf("encode llm request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("create llm request: %w", err)
	}
	if p.apiKey != "" {
		request.Header.Set("Authorization", "Bearer "+p.apiKey)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := p.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("send llm request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("llm request returned HTTP status %d", response.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 4<<20)).Decode(&result); err != nil {
		return "", fmt.Errorf("decode llm response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", errors.New("llm response contains no choices")
	}
	return result.Choices[0].Message.Content, nil
}
