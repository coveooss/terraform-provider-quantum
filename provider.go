package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"quantum_file":       dataSourceQuantumFile(),
			"quantum_query_json": dataSourceQuantumQueryJSON(),
			"quantum_list_files": dataSourceQuantumListFiles(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"quantum_password": resourceQuantumPassword(),
		},
	}
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}
