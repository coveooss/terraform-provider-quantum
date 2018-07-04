package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tidwall/gjson"
)

func dataSourceQuantumQueryJSON() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQuantumQueryJSONRead,

		Schema: map[string]*schema.Schema{
			"json": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"query": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"result": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceQuantumQueryJSONRead(d *schema.ResourceData, m interface{}) error {
	json := d.Get("json").(string)
	query := d.Get("query").(string)
	queryResult := gjson.Get(json, query)

	d.SetId(fmt.Sprintf("%d-%d", hashcode.String(json), hashcode.String(query)))
	d.Set("result", queryResult.String())

	return nil
}
