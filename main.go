package main

import (
	"github.com/echohello-dev/terraform-provider-openrouter/openrouter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: openrouter.Provider,
	}

	plugin.Serve(opts)
}
