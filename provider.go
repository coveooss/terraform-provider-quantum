package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return &schema.Provider{
				DataSourcesMap: map[string]*schema.Resource{
					"quantum_list_files":    dataSourceQuantumListFiles(),
				},
			}
		},
	})
}
