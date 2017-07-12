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
					"quantum_list_files": dataSourceQuantumListFiles(),
					// "quantum_password":   dataSourceQuantumPassword(),
				},
				ResourcesMap: map[string]*schema.Resource{
					"quantum_password": resourceQuantumPassword(),
				},

				// ConfigureFunc: configureQuantum,
			}

		},
	})
}

// func configureQuantum(providerSettings *schema.ResourceData) (interface{}, error) {

// 	return &QuantumMeta{
// 		passwords: make(map[string]string),
// 	}, nil
// }

// // QuantumMeta contains existing generated password from state file
// type QuantumMeta struct {
// 	passwords map[string]string
// }
