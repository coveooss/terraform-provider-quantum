package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
         content     = "This is some content"
         filename    = "quantum_file"
      }`,
		},
		{
			"quantum_file",
			"This is some sensitive content",
			`data "quantum_file" "file" {
         sensitive_content     = "This is some sensitive content"
         filename    = "quantum_file"
      }`,
		},
	}

	for _, tt := range cases {
		r.UnitTest(t, r.TestCase{
			Providers: testProviders,
			Steps: []r.TestStep{
				{
					Config: tt.config,
					Check: func(s *terraform.State) error {
						content, err := ioutil.ReadFile(tt.path)
						if err != nil {
							return fmt.Errorf("config:\n%s\n,got: %s\n", tt.config, err)
						}
						if string(content) != tt.content {
							return fmt.Errorf("config:\n%s\ngot:\n%s\nwant:\n%s\n", tt.config, content, tt.content)
						}
						return nil
					},
				},
			},
		})
	}

	os.Remove("quantum_file")
}
