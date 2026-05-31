package openrouter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: modelsRead,
		Schema: map[string]*schema.Schema{
			"output_modality": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter models by output modality. Valid values: text, image, audio, embeddings, all.",
				ValidateFunc: validation.StringInSlice([]string{"text", "image", "audio", "embeddings", "all"}, false),
			},
			"models": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier for the model.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the model.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the model.",
						},
						"context_length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum context length in tokens.",
						},
						"created": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp of when the model was created.",
						},
						"pricing": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeFloat},
							Description: "Pricing information for the model (input price per million tokens).",
						},
						"top_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Top provider for this model.",
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Model architecture.",
						},
						"supported_parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of supported parameters for this model.",
						},
						"knowledge_cutoff": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date up to which the model was trained on data.",
						},
						"features": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Feature flags for the model.",
						},
					},
				},
			},
		},
	}
}

func modelsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	outputModality := d.Get("output_modality").(string)

	models, diags := client.ListModels(ctx, outputModality)
	if diags.HasError() {
		return diags
	}

	if models == nil {
		d.SetId("openrouter-models")
		if err := d.Set("models", []interface{}{}); err != nil {
			return diag.Errorf("Failed to set models: %s", err)
		}
		return nil
	}

	modelsList := make([]map[string]interface{}, 0, len(models))
	for _, m := range models {
		modelMap := map[string]interface{}{
			"id":                   m.ID,
			"name":                 m.Name,
			"description":          m.Description,
			"context_length":       m.ContextLength,
			"created":              m.Created,
			"top_provider":         m.TopProvider.SelectedMode,
			"architecture":         m.Architecture.ModelClass,
			"knowledge_cutoff":      m.KnowledgeCutoff,
			"features":             m.Features,
			"supported_parameters":  m.SupportedParams,
		}

		pricing := make(map[string]interface{})
		for tier, price := range m.Pricing {
			pricing[tier] = price.Input
		}
		modelMap["pricing"] = pricing

		modelsList = append(modelsList, modelMap)
	}

	if err := d.Set("models", modelsList); err != nil {
		return diag.Errorf("Failed to set models: %s", err)
	}

	d.SetId(fmt.Sprintf("openrouter-models-%d", len(models)))

	return nil
}
