package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tidwall/gjson"
)

func dataSourceQuantumQueryJSON() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Terraform 0.12 supports nested maps. Ex: jsondecode(var.my_variable).attribute1.attribute2",
		Read:               dataSourceQuantumQueryJSONRead,

		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result_list": {
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Type:     schema.TypeList,
			},
			"result_map": {
				Computed: true,
				Type:     schema.TypeMap,
			},
			"result": {
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

	if queryResult.IsArray() {
		resultList := []string{}
		for _, value := range queryResult.Array() {
			resultList = append(resultList, value.String())
		}
		d.Set("result_list", resultList)
	}
	if queryResult.IsObject() {
		resultMap := map[string]string{}
		for key, value := range queryResult.Map() {
			resultMap[key] = value.String()
		}
		d.Set("result_map", resultMap)
	}

	return nil
}
