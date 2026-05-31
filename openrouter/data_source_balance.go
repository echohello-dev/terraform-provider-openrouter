package openrouter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBalance() *schema.Resource {
	return &schema.Resource{
		ReadContext: balanceRead,
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Label for the API key.",
			},
			"limit": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Credit limit for the key, or null if unlimited.",
			},
			"limit_remaining": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Remaining credits for the key, or null if unlimited.",
			},
			"usage": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Number of credits used (all time).",
			},
			"usage_daily": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Number of credits used (current UTC day).",
			},
			"usage_weekly": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Number of credits used (current UTC week, starting Monday).",
			},
			"usage_monthly": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Number of credits used (current UTC month).",
			},
			"is_free_tier": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user has paid for credits before.",
			},
		},
	}
}

func balanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	keyData, err := client.GetBalance(ctx)
	if err != nil {
		return err
	}

	d.SetId("openrouter-balance")
	if err := d.Set("label", keyData.Label); err != nil {
		return diag.Errorf("Failed to set label: %s", err)
	}
	if err := d.Set("limit", keyData.Limit); err != nil {
		return diag.Errorf("Failed to set limit: %s", err)
	}
	if err := d.Set("limit_remaining", keyData.LimitRemaining); err != nil {
		return diag.Errorf("Failed to set limit_remaining: %s", err)
	}
	if err := d.Set("usage", keyData.Usage); err != nil {
		return diag.Errorf("Failed to set usage: %s", err)
	}
	if err := d.Set("usage_daily", keyData.UsageDaily); err != nil {
		return diag.Errorf("Failed to set usage_daily: %s", err)
	}
	if err := d.Set("usage_weekly", keyData.UsageWeekly); err != nil {
		return diag.Errorf("Failed to set usage_weekly: %s", err)
	}
	if err := d.Set("usage_monthly", keyData.UsageMonthly); err != nil {
		return diag.Errorf("Failed to set usage_monthly: %s", err)
	}
	if err := d.Set("is_free_tier", keyData.IsFreeTier); err != nil {
		return diag.Errorf("Failed to set is_free_tier: %s", err)
	}

	return nil
}
