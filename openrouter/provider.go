package openrouter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENROUTER_API_KEY", nil),
				Description: "OpenRouter API key. Can also be set via OPENROUTER_API_KEY environment variable.",
			},
			"referer": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENROUTER_REFERER", ""),
				Description: "Site URL for OpenRouter rankings. Can also be set via OPENROUTER_REFERER environment variable.",
			},
			"app_title": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENROUTER_APP_TITLE", ""),
				Description: "App display name for OpenRouter dashboard. Can also be set via OPENROUTER_APP_TITLE environment variable.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"openrouter_models":     dataSourceModels(),
			"openrouter_model":      dataSourceModel(),
			"openrouter_balance":    dataSourceBalance(),
			"openrouter_generation": dataSourceGeneration(),
			"openrouter_credits":    dataSourceCredits(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"openrouter_chat_completion": resourceChatCompletion(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	referer := d.Get("referer").(string)
	appTitle := d.Get("app_title").(string)

	client := NewClient(apiKey, referer, appTitle)

	return client, nil
}
