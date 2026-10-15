package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	openAIAPIKeyEnv = "OPENAI_API_KEY"
	openAIModel     = "gpt-4o-mini"
	openAIEndpoint  = "https://api.openai.com/v1/chat/completions"
)

type OpenAIProvider struct {
	apiKey string
	client *http.Client
}

func NewOpenAIProvider() (*OpenAIProvider, error) {
	apiKey := strings.TrimSpace(os.Getenv(openAIAPIKeyEnv))
	if apiKey == "" {
		return nil, fmt.Errorf("%s environment variable is required for llm steps", openAIAPIKeyEnv)
	}
	return &OpenAIProvider{apiKey: apiKey, client: http.DefaultClient}, nil
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt string) (string, error) {
	requestBody, err := json.Marshal(struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}{
		Model: openAIModel,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return "", fmt.Errorf("encode OpenAI request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("create OpenAI request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+p.apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := p.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("send OpenAI request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("OpenAI request returned HTTP status %d", response.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 4<<20)).Decode(&result); err != nil {
		return "", fmt.Errorf("decode OpenAI response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", errors.New("OpenAI response contains no choices")
	}
	return result.Choices[0].Message.Content, nil
}
