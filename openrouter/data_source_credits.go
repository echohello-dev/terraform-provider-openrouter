package openrouter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCredits() *schema.Resource {
	return &schema.Resource{
		ReadContext: creditsRead,
		Schema: map[string]*schema.Schema{
			"total_credits": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Total credits purchased.",
			},
			"total_usage": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Total credits used.",
			},
		},
	}
}

func creditsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	credits, err := client.GetCredits(ctx)
	if err != nil {
		return err
	}

	d.SetId("openrouter-credits")
	if err := d.Set("total_credits", credits.TotalCredits); err != nil {
		return diag.Errorf("Failed to set total_credits: %s", err)
	}
	if err := d.Set("total_usage", credits.TotalUsage); err != nil {
		return diag.Errorf("Failed to set total_usage: %s", err)
	}

	return nil
}
