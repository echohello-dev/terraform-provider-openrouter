package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const (
	BaseURL            = "https://openrouter.ai/api/v1"
	DefaultTimeout     = 120 * time.Second
	MaxRetries         = 3
	RetryBaseDelay     = 1 * time.Second
)

type Client struct {
	APIKey      string
	Referer     string
	AppTitle    string
	HTTPClient  *http.Client
}

type Model struct {
	ID                 string                `json:"id"`
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	Created            int64                 `json:"created"`
	ContextLength      int64                 `json:"context_length"`
	Pricing            map[string]PriceTier  `json:"pricing"`
	TopProvider        TopProviderInfo       `json:"top_provider"`
	Architecture       ModelArchitecture     `json:"architecture"`
	Features           []string              `json:"features,omitempty"`
	SupportedParams    []string              `json:"supported_parameters,omitempty"`
	KnowledgeCutoff   string                `json:"knowledge_cutoff,omitempty"`
	CanonicalSlug      string                `json:"canonical_slug,omitempty"`
	HuggingFaceID      string                `json:"hugging_face_id,omitempty"`
	ExpirationDate     string                `json:"expiration_date,omitempty"`
}

type PriceTier struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

type TopProviderInfo struct {
	NumModes     int64    `json:"num_modes,omitempty"`
	SelectedMode string   `json:"selected_mode,omitempty"`
	Modes        []string `json:"modes,omitempty"`
}

type ModelArchitecture struct {
	BaseModelClasses []string `json:"base_model_classes,omitempty"`
	ModelClass       string   `json:"model_class,omitempty"`
}

type ModelsResponse struct {
	Data []Model `json:"data"`
}

type BalanceResponse struct {
	Data BalanceData `json:"data"`
}

type BalanceData struct {
	TotalCredits float64 `json:"totalCredits,omitempty"`
	UsedCredits  float64 `json:"usedCredits,omitempty"`
}

type KeyResponse struct {
	Data KeyData `json:"data"`
}

type KeyData struct {
	Label            string   `json:"label"`
	Limit            *float64 `json:"limit,omitempty"`
	LimitRemaining   *float64 `json:"limit_remaining,omitempty"`
	Usage            float64  `json:"usage"`
	UsageDaily       float64  `json:"usage_daily"`
	UsageWeekly      float64  `json:"usage_weekly"`
	UsageMonthly     float64  `json:"usage_monthly"`
	IsFreeTier       bool     `json:"is_free_tier"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      int                    `json:"code,omitempty"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type APIError struct {
	StatusCode int
	Code       int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

func NewClient(apiKey, referer, appTitle string) *Client {
	return &Client{
		APIKey:   apiKey,
		Referer:  referer,
		AppTitle: appTitle,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, *http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	if c.Referer != "" {
		req.Header.Set("HTTP-Referer", c.Referer)
	}
	if c.AppTitle != "" {
		req.Header.Set("X-OpenRouter-Title", c.AppTitle)
	}

	var lastErr error
	for attempt := 0; attempt < MaxRetries; attempt++ {
		if attempt > 0 {
			delay := RetryBaseDelay * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to execute request: %w", err)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		if resp.StatusCode == 429 {
			lastErr = fmt.Errorf("rate limited")
			continue
		}

		if resp.StatusCode >= 400 {
			var errResp ErrorResponse
			if unmarshalErr := json.Unmarshal(respBody, &errResp); unmarshalErr == nil {
				return nil, resp, &APIError{
					StatusCode: resp.StatusCode,
					Code:       errResp.Error.Code,
					Message:    errResp.Error.Message,
				}
			}
			return nil, resp, &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}

		return respBody, resp, nil
	}

	return nil, nil, lastErr
}

func (c *Client) ListModels(ctx context.Context, outputModality string) ([]Model, diag.Diagnostics) {
	endpoint := "/models"
	if outputModality != "" {
		endpoint += "?output_modalities=" + outputModality
	}

	body, _, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, diag.Errorf("Failed to list models: %s", err)
	}

	var resp ModelsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, diag.Errorf("Failed to unmarshal response: %s", err)
	}

	return resp.Data, nil
}

func (c *Client) GetModel(ctx context.Context, modelID string) (*Model, diag.Diagnostics) {
	body, _, err := c.doRequest(ctx, "GET", "/models", nil)
	if err != nil {
		return nil, diag.Errorf("Failed to get model: %s", err)
	}

	var resp ModelsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, diag.Errorf("Failed to unmarshal response: %s", err)
	}

	for _, m := range resp.Data {
		if m.ID == modelID {
			return &m, nil
		}
	}

	return nil, diag.Errorf("Model not found: %s", modelID)
}

func (c *Client) GetBalance(ctx context.Context) (*KeyData, diag.Diagnostics) {
	body, _, err := c.doRequest(ctx, "GET", "/key", nil)
	if err != nil {
		return nil, diag.Errorf("Failed to get balance: %s", err)
	}

	var keyResp KeyResponse
	if err := json.Unmarshal(body, &keyResp); err != nil {
		return nil, diag.Errorf("Failed to unmarshal response: %s", err)
	}

	return &keyResp.Data, nil
}

func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, diag.Diagnostics) {
	body, _, err := c.doRequest(ctx, "POST", "/chat/completions", req)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			switch apiErr.StatusCode {
			case 401:
				return nil, diag.Errorf("Invalid API key. Please check your OpenRouter API key.")
			case 402:
				return nil, diag.Errorf("Insufficient credits. Please add credits to your OpenRouter account.")
			case 429:
				return nil, diag.Errorf("Rate limited. Please try again later.")
			default:
				return nil, diag.Errorf("API error: %s", apiErr.Message)
			}
		}
		return nil, diag.Errorf("Failed to create chat completion: %s", err)
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, diag.Errorf("Failed to unmarshal response: %s", err)
	}

	if len(chatResp.Choices) == 0 && chatResp.Usage.TotalTokens == 0 {
		return nil, diag.Errorf("Empty response from API")
	}

	return &chatResp, nil
}
