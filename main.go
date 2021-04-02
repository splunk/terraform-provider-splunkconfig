package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-splunkconfig/internal/provider"
)

const (
	version = "dev"
)

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: provider.New(version)}

	plugin.Serve(opts)
}
