package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceQuantumFile() *schema.Resource {
	return &schema.Resource{
		Read: resourceLocalFileRead,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"sensitive_content"},
			},
			"sensitive_content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{"content"},
			},
			"filename": {
				Type:        schema.TypeString,
				Description: "Path to the output file",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceLocalFileRead(d *schema.ResourceData, _ interface{}) error {
	content := resourceLocalFileContent(d)
	destination := d.Get("filename").(string)

	destinationDir := path.Dir(destination)
	if _, err := os.Stat(destinationDir); err != nil {
		if err := os.MkdirAll(destinationDir, 0777); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(destination, []byte(content), 0777); err != nil {
		return err
	}

	checksum := sha1.Sum([]byte(content))
	d.SetId(hex.EncodeToString(checksum[:]))

	return nil
}

func resourceLocalFileContent(d *schema.ResourceData) string {
	content := d.Get("content")
	sensitiveContent, sensitiveSpecified := d.GetOk("sensitive_content")
	useContent := content.(string)
	if sensitiveSpecified {
		useContent = sensitiveContent.(string)
	}

	return useContent
}
