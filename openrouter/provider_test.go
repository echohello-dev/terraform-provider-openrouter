package openrouter

import (
	"encoding/json"
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("Provider failed validation: %s", err)
	}
}

func TestProviderSchema(t *testing.T) {
	provider := Provider()

	schema := provider.Schema

	if _, ok := schema["api_key"]; !ok {
		t.Error("Provider schema should have api_key field")
	}

	if _, ok := schema["referer"]; !ok {
		t.Error("Provider schema should have referer field")
	}

	if _, ok := schema["app_title"]; !ok {
		t.Error("Provider schema should have app_title field")
	}
}

func TestDataSourceModelsSchema(t *testing.T) {
	ds := dataSourceModels()

	if ds == nil {
		t.Fatal("dataSourceModels() returned nil")
	}

	schema := ds.Schema

	if _, ok := schema["output_modality"]; !ok {
		t.Error("Models data source should have output_modality field")
	}

	if _, ok := schema["models"]; !ok {
		t.Error("Models data source should have models field")
	}
}

func TestDataSourceBalanceSchema(t *testing.T) {
	ds := dataSourceBalance()

	if ds == nil {
		t.Fatal("dataSourceBalance() returned nil")
	}

	schema := ds.Schema

	expectedFields := []string{
		"label",
		"limit",
		"limit_remaining",
		"usage",
		"usage_daily",
		"usage_weekly",
		"usage_monthly",
		"is_free_tier",
	}

	for _, field := range expectedFields {
		if _, ok := schema[field]; !ok {
			t.Errorf("Balance data source should have %s field", field)
		}
	}
}

func TestResourceChatCompletionSchema(t *testing.T) {
	r := resourceChatCompletion()

	if r == nil {
		t.Fatal("resourceChatCompletion() returned nil")
	}

	schema := r.Schema

	requiredFields := []string{
		"model", "messages", "max_tokens", "temperature", "top_p",
		"seed", "stop", "frequency_penalty", "presence_penalty",
		"response_format", "stream", "user", "session_id",
		"logprobs", "top_logprobs", "content",
	}

	for _, field := range requiredFields {
		if _, ok := schema[field]; !ok {
			t.Errorf("Chat completion resource should have %s field", field)
		}
	}
}

func TestDataSourceGenerationSchema(t *testing.T) {
	ds := dataSourceGeneration()

	if ds == nil {
		t.Fatal("dataSourceGeneration() returned nil")
	}

	schema := ds.Schema

	expectedFields := []string{
		"id", "upstream_id", "total_cost", "cache_discount",
		"upstream_inference_cost", "created_at", "model", "provider_name",
		"latency", "generation_time", "finish_reason", "tokens_prompt",
		"tokens_completion", "native_tokens_prompt", "native_tokens_completion",
		"native_tokens_reasoning", "native_tokens_cached", "num_media_prompt",
		"num_media_completion", "origin", "usage", "is_byok",
		"native_finish_reason", "api_type", "router",
	}

	for _, field := range expectedFields {
		if _, ok := schema[field]; !ok {
			t.Errorf("Generation data source should have %s field", field)
		}
	}
}

func TestDataSourceCreditsSchema(t *testing.T) {
	ds := dataSourceCredits()

	if ds == nil {
		t.Fatal("dataSourceCredits() returned nil")
	}

	schema := ds.Schema

	if _, ok := schema["total_credits"]; !ok {
		t.Error("Credits data source should have total_credits field")
	}
	if _, ok := schema["total_usage"]; !ok {
		t.Error("Credits data source should have total_usage field")
	}
}

func TestGenerationResponseParsing(t *testing.T) {
	jsonResp := `{
		"data": {
			"id": "gen-abc123",
			"upstream_id": "upstream-xyz",
			"total_cost": 0.0012,
			"cache_discount": 0.0001,
			"upstream_inference_cost": 0.0011,
			"created_at": "2024-01-15T10:30:00Z",
			"model": "openai/gpt-4",
			"provider_name": "OpenAI",
			"latency": 450,
			"generation_time": 380,
			"finish_reason": "stop",
			"tokens_prompt": 10,
			"tokens_completion": 25,
			"native_tokens_prompt": 10,
			"native_tokens_completion": 25,
			"native_tokens_reasoning": 0,
			"native_tokens_cached": 0,
			"num_media_prompt": 0,
			"num_media_completion": 0,
			"origin": "https://example.com",
			"usage": 0.0012,
			"is_byok": false,
			"native_finish_reason": "stop",
			"api_type": "completions",
			"router": "openrouter/auto"
		}
	}`

	var resp GenerationResponse
	if err := json.Unmarshal([]byte(jsonResp), &resp); err != nil {
		t.Fatalf("Failed to unmarshal generation response: %s", err)
	}

	if resp.Data.ID != "gen-abc123" {
		t.Errorf("Expected ID to be 'gen-abc123', got '%s'", resp.Data.ID)
	}

	if resp.Data.TotalCost == nil || *resp.Data.TotalCost != 0.0012 {
		t.Errorf("Expected total_cost to be 0.0012")
	}

	if resp.Data.ProviderName == nil || *resp.Data.ProviderName != "OpenAI" {
		t.Errorf("Expected provider_name to be 'OpenAI'")
	}
}

func TestCreditsResponseParsing(t *testing.T) {
	jsonResp := `{"data":{"total_credits":100.0,"total_usage":23.45}}`

	var resp CreditsResponse
	if err := json.Unmarshal([]byte(jsonResp), &resp); err != nil {
		t.Fatalf("Failed to unmarshal credits response: %s", err)
	}

	if resp.Data.TotalCredits != 100.0 {
		t.Errorf("Expected total_credits to be 100.0, got %f", resp.Data.TotalCredits)
	}

	if resp.Data.TotalUsage != 23.45 {
		t.Errorf("Expected total_usage to be 23.45, got %f", resp.Data.TotalUsage)
	}
}

func TestClientNewClient(t *testing.T) {
	client := NewClient("test-key", "https://example.com", "TestApp")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.APIKey != "test-key" {
		t.Errorf("Expected APIKey to be 'test-key', got '%s'", client.APIKey)
	}

	if client.Referer != "https://example.com" {
		t.Errorf("Expected Referer to be 'https://example.com', got '%s'", client.Referer)
	}

	if client.AppTitle != "TestApp" {
		t.Errorf("Expected AppTitle to be 'TestApp', got '%s'", client.AppTitle)
	}

	if client.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}
}

func TestChatCompletionResponseParsing(t *testing.T) {
	jsonResp := `{
		"id": "chatcmpl-123",
		"object": "chat.completion",
		"created": 1677652288,
		"model": "gpt-4",
		"choices": [{
			"index": 0,
			"message": {
				"role": "assistant",
				"content": "Hello!"
			},
			"finish_reason": "stop"
		}],
		"usage": {
			"prompt_tokens": 10,
			"completion_tokens": 5,
			"total_tokens": 15
		}
	}`

	var resp ChatCompletionResponse
	if err := json.Unmarshal([]byte(jsonResp), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}

	if resp.ID != "chatcmpl-123" {
		t.Errorf("Expected ID to be 'chatcmpl-123', got '%s'", resp.ID)
	}

	if resp.Model != "gpt-4" {
		t.Errorf("Expected Model to be 'gpt-4', got '%s'", resp.Model)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Hello!" {
		t.Errorf("Expected content to be 'Hello!', got '%s'", resp.Choices[0].Message.Content)
	}

	if resp.Usage.TotalTokens != 15 {
		t.Errorf("Expected total_tokens to be 15, got %d", resp.Usage.TotalTokens)
	}
}
