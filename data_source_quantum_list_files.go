package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceQuantumListFiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQuantumListFilesRead,

		Schema: map[string]*schema.Schema{
			"folders": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"patterns": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_folder": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"recursive": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"files": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceQuantumListFilesRead(d *schema.ResourceData, m interface{}) error {
	folders := interfaceToString(d.Get("folders").([]interface{}))
	if len(folders) == 0 {
		folders = []string{"."}
	}

	patterns := interfaceToString(d.Get("patterns").([]interface{}))
	if len(patterns) == 0 {
		patterns = []string{"*"}
	}

	includeFolder := d.Get("include_folder").(bool)
	recursive := d.Get("recursive").(bool)

	var result []string
	for _, folder := range folders {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			for _, pattern := range patterns {
				matched, err := filepath.Match(pattern, filepath.Base(path))
				if err != nil {
					return err
				}
				if matched {
					addedFile := filepath.ToSlash(path)
					if !includeFolder {
						if addedFile == folder {
							continue
						}

						addedFile = strings.TrimPrefix(addedFile, folder+"/")
					}
					result = append(result, addedFile)
					break
				}
			}

			if os.FileInfo.IsDir(info) && !recursive && path != folder {
				return filepath.SkipDir
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	d.Set("files", result)
	d.SetId("-")

	return nil
}
