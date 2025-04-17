package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestFile_basic(t *testing.T) {
	var cases = []struct {
		path    string
		content string
		config  string
	}{
		{
			"quantum_file",
			"This is some content",
			`data "quantum_file" "file" {
				content  = "This is some content"
				filename = "quantum_file"
			}`,
		},
		{
			"quantum_file",
			"This is some sensitive content",
			`data "quantum_file" "file" {
				sensitive_content = "This is some sensitive content"
				filename          = "quantum_file"
			}`,
		},
	}

	for _, tt := range cases {
		var testProviders = map[string]*schema.Provider{
			"quantum": Provider(),
		}

		resource.UnitTest(t, resource.TestCase{
			Providers: testProviders, // assumes testProviders is defined as: map[string]*schema.Provider
			Steps: []resource.TestStep{
				{
					Config: tt.config,
					Check: resource.TestCheckFunc(func(s *terraform.State) error {
						content, err := os.ReadFile(tt.path)
						if err != nil {
							return fmt.Errorf("config:\n%s\n, got error: %s\n", tt.config, err)
						}
						if string(content) != tt.content {
							return fmt.Errorf("config:\n%s\ngot:\n%s\nwant:\n%s\n", tt.config, content, tt.content)
						}
						return nil
					}),
				},
			},
		})
	}

	os.Remove("quantum_file")
}
