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

	if _, ok := schema["model"]; !ok {
		t.Error("Chat completion resource should have model field")
	}

	if _, ok := schema["messages"]; !ok {
		t.Error("Chat completion resource should have messages field")
	}

	if _, ok := schema["max_tokens"]; !ok {
		t.Error("Chat completion resource should have max_tokens field")
	}

	if _, ok := schema["temperature"]; !ok {
		t.Error("Chat completion resource should have temperature field")
	}

	if _, ok := schema["content"]; !ok {
		t.Error("Chat completion resource should have computed content field")
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
