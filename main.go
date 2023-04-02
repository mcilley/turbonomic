package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic"
)

func main() {
	// opts contains the configurations to serve the Turbonomic plugin.
	opts := plugin.ServeOpts{
		ProviderFunc: turbonomic.Provider,
	}
	// Serves the Turbonomic plugin in the defined configurations.
	plugin.Serve(&opts)
}
