package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceQuantumListFiles() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Please use `fileset` instead: https://www.terraform.io/docs/configuration/functions/fileset.html",
		Read:               dataSourceQuantumListFilesRead,

		Schema: map[string]*schema.Schema{
			"folders": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"patterns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_folder": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"recursive": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"files": {
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
		folderInfo, err := os.Stat(folder)
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", folder)
		}
		if !folderInfo.IsDir() {
			return fmt.Errorf("%s is not a dir", folder)
		}

		err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
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
