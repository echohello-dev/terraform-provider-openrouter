package openrouter

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceChatCompletion() *schema.Resource {
	return &schema.Resource{
		CreateContext: chatCompletionCreate,
		ReadContext:   chatCompletionRead,
		DeleteContext: chatCompletionDelete,
		Schema: map[string]*schema.Schema{
			"model": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The model to use for the completion. Example: openai/gpt-4.",
				ValidateFunc: validation.NoZeroValues,
			},
			"messages": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Role of the message. Must be one of: system, user, assistant.",
							ValidateFunc: validation.StringInSlice([]string{"system", "user", "assistant", "tool"}, false),
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Content of the message.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the sender (optional).",
						},
					},
				},
			},
			"max_tokens": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Description:  "Maximum number of tokens to generate.",
				ValidateFunc: validation.IntAtLeast(1),
			},
		"temperature": {
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Description:  "Sampling temperature (0-2). Higher values make output more random. If unset, the API default (1.0) is used.",
			ValidateFunc: validation.FloatBetween(0, 2),
		},
		"top_p": {
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Description:  "Nucleus sampling parameter. Higher values make output more focused.",
			ValidateFunc: validation.FloatBetween(0, 1),
		},
		"seed": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Description:  "Seed for deterministic sampling. If set, repeated requests with the same parameters should return similar results.",
		},
		"stop": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    4,
			Description: "Up to 4 sequences where the API will stop generating further tokens.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"frequency_penalty": {
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Description:  "Penalty for repeating tokens (-2.0 to 2.0). Positive values reduce repetition.",
			ValidateFunc: validation.FloatBetween(-2, 2),
		},
		"presence_penalty": {
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Description:  "Penalty for introducing new tokens (-2.0 to 2.0). Positive values encourage new topics.",
			ValidateFunc: validation.FloatBetween(-2, 2),
		},
		"response_format": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "Format for the response: text, json_object, or json_schema.",
			ValidateFunc: validation.StringInSlice([]string{"text", "json_object", "json_schema"}, false),
		},
		"stream": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Whether to stream the response back. Terraform cannot consume streams; this is for API passthrough only.",
		},
		"user": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "A unique identifier representing your end-user.",
		},
		"session_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "A stable session identifier (max 128 chars).",
			ValidateFunc: validation.StringLenBetween(1, 128),
		},
		"logprobs": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Whether to return log probabilities of the output tokens.",
		},
		"top_logprobs": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Description:  "Number of top log probabilities to return per output token (0-20). Only used if logprobs is true.",
			ValidateFunc: validation.IntBetween(0, 20),
		},
			"response_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for this completion.",
			},
			"response_object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object type (e.g., chat.completion).",
			},
			"response_created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp of when the response was created.",
			},
			"response_model": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Actual model used for the completion.",
			},
			"finish_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason the completion finished (stop, length, etc.).",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated completion content.",
			},
			"prompt_tokens": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tokens in the prompt.",
			},
			"completion_tokens": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tokens in the completion.",
			},
			"total_tokens": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of tokens used.",
			},
		},

	}
}

func chatCompletionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	model := d.Get("model").(string)
	maxTokens := d.Get("max_tokens").(int)

	messagesRaw := d.Get("messages").([]interface{})
	messages := make([]ChatMessage, 0, len(messagesRaw))
	for _, m := range messagesRaw {
		msgMap := m.(map[string]interface{})
		content := strings.TrimSpace(msgMap["content"].(string))
		if content == "" {
			continue
		}
		msg := ChatMessage{
			Role:    strings.ToLower(msgMap["role"].(string)),
			Content: content,
		}
		if name, ok := msgMap["name"].(string); ok && name != "" {
			msg.Name = name
		}
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return diag.Errorf("At least one message with non-empty content is required")
	}

	req := ChatCompletionRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: maxTokens,
		Stream:    d.Get("stream").(bool),
		Logprobs:  d.Get("logprobs").(bool),
		User:      d.Get("user").(string),
		SessionID: d.Get("session_id").(string),
	}
	if v, ok := d.GetOk("seed"); ok {
		s := v.(int)
		req.Seed = &s
	}
	if v, ok := d.GetOk("stop"); ok {
		stopRaw := v.([]interface{})
		stops := make([]string, 0, len(stopRaw))
		for _, s := range stopRaw {
			if str := s.(string); str != "" {
				stops = append(stops, str)
			}
		}
		req.Stop = stops
	}
	if tempVal := d.GetRawConfig().GetAttr("temperature"); !tempVal.IsNull() {
		t, _ := tempVal.AsBigFloat().Float64()
		req.Temperature = &t
	}
	if topPVal := d.GetRawConfig().GetAttr("top_p"); !topPVal.IsNull() {
		p, _ := topPVal.AsBigFloat().Float64()
		req.TopP = &p
	}
	if freqVal := d.GetRawConfig().GetAttr("frequency_penalty"); !freqVal.IsNull() {
		fp, _ := freqVal.AsBigFloat().Float64()
		req.FrequencyPenalty = &fp
	}
	if presVal := d.GetRawConfig().GetAttr("presence_penalty"); !presVal.IsNull() {
		pp, _ := presVal.AsBigFloat().Float64()
		req.PresencePenalty = &pp
	}
	if v, ok := d.GetOk("top_logprobs"); ok {
		tlp := v.(int)
		req.TopLogprobs = &tlp
	}
	if v, ok := d.GetOk("response_format"); ok {
		rf := v.(string)
		switch rf {
		case "json_object":
			req.ResponseFormat = ResponseFormatJSONObject{Type: "json_object"}
		default:
			req.ResponseFormat = ResponseFormatText{Type: "text"}
		}
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		if err := d.Set("finish_reason", choice.FinishReason); err != nil {
			return diag.Errorf("Failed to set finish_reason: %s", err)
		}
		if err := d.Set("content", choice.Message.Content); err != nil {
			return diag.Errorf("Failed to set content: %s", err)
		}
	}

	if err := d.Set("response_id", resp.ID); err != nil {
		return diag.Errorf("Failed to set response_id: %s", err)
	}
	if err := d.Set("response_object", resp.Object); err != nil {
		return diag.Errorf("Failed to set response_object: %s", err)
	}
	if err := d.Set("response_created", resp.Created); err != nil {
		return diag.Errorf("Failed to set response_created: %s", err)
	}
	if err := d.Set("response_model", resp.Model); err != nil {
		return diag.Errorf("Failed to set response_model: %s", err)
	}
	if err := d.Set("prompt_tokens", resp.Usage.PromptTokens); err != nil {
		return diag.Errorf("Failed to set prompt_tokens: %s", err)
	}
	if err := d.Set("completion_tokens", resp.Usage.CompletionTokens); err != nil {
		return diag.Errorf("Failed to set completion_tokens: %s", err)
	}
	if err := d.Set("total_tokens", resp.Usage.TotalTokens); err != nil {
		return diag.Errorf("Failed to set total_tokens: %s", err)
	}

	return nil
}

func chatCompletionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func chatCompletionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
