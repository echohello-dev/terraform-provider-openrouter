package openrouter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGeneration() *schema.Resource {
	return &schema.Resource{
		ReadContext: generationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The generation ID (returned as response_id from openrouter_chat_completion).",
			},
			"upstream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provider's upstream generation ID.",
			},
			"total_cost": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Total cost of the generation in USD.",
			},
			"cache_discount": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Discount from caching in USD.",
			},
			"upstream_inference_cost": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Cost charged by the upstream provider in USD.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ISO 8601 timestamp when the generation was created.",
			},
			"model": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Model used for the generation.",
			},
			"provider_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the provider that served the request.",
			},
			"latency": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Total latency in milliseconds.",
			},
			"generation_time": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Generation time in milliseconds.",
			},
			"finish_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason the generation finished.",
			},
			"tokens_prompt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of prompt tokens.",
			},
			"tokens_completion": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of completion tokens.",
			},
			"native_tokens_prompt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Provider-reported prompt tokens.",
			},
			"native_tokens_completion": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Provider-reported completion tokens.",
			},
			"native_tokens_reasoning": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Provider-reported reasoning tokens.",
			},
			"native_tokens_cached": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Provider-reported cached tokens.",
			},
			"num_media_prompt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of media items in the prompt.",
			},
			"num_media_completion": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of media items in the completion.",
			},
			"origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Origin URL of the request.",
			},
			"usage": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Usage cost in USD.",
			},
			"is_byok": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the request used a Bring Your Own Key configuration.",
			},
			"native_finish_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provider's native finish reason.",
			},
			"api_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API type used (e.g. completions, embeddings).",
			},
			"router": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Router used (e.g. openrouter/auto).",
			},
		},
	}
}

func generationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	genID := d.Get("id").(string)
	if genID == "" {
		return diag.Errorf("id is required")
	}

	gen, err := client.GetGeneration(ctx, genID)
	if err != nil {
		return err
	}

	d.SetId(gen.ID)

	setString := func(key string, val *string) {
		if val != nil {
			_ = d.Set(key, *val)
		}
	}
	setFloat := func(key string, val *float64) {
		if val != nil {
			_ = d.Set(key, *val)
		}
	}
	setInt := func(key string, val *int) {
		if val != nil {
			_ = d.Set(key, *val)
		}
	}

	setString("upstream_id", gen.UpstreamID)
	setFloat("total_cost", gen.TotalCost)
	setFloat("cache_discount", gen.CacheDiscount)
	setFloat("upstream_inference_cost", gen.UpstreamInferenceCost)
	_ = d.Set("created_at", gen.CreatedAt)
	_ = d.Set("model", gen.Model)
	setString("provider_name", gen.ProviderName)
	setFloat("latency", gen.Latency)
	setFloat("generation_time", gen.GenerationTime)
	setString("finish_reason", gen.FinishReason)
	setInt("tokens_prompt", gen.TokensPrompt)
	setInt("tokens_completion", gen.TokensCompletion)
	setInt("native_tokens_prompt", gen.NativeTokensPrompt)
	setInt("native_tokens_completion", gen.NativeTokensCompletion)
	setInt("native_tokens_reasoning", gen.NativeTokensReasoning)
	setInt("native_tokens_cached", gen.NativeTokensCached)
	setInt("num_media_prompt", gen.NumMediaPrompt)
	setInt("num_media_completion", gen.NumMediaCompletion)
	_ = d.Set("origin", gen.Origin)
	_ = d.Set("usage", gen.Usage)
	_ = d.Set("is_byok", gen.IsByok)
	setString("native_finish_reason", gen.NativeFinishReason)
	setString("api_type", gen.APIType)
	setString("router", gen.Router)

	return nil
}
