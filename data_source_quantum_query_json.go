package main

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
)

func computeHash(s string) int {
	sum := sha1.Sum([]byte(s))
	return int(binary.BigEndian.Uint32(sum[0:4]))
}

func dataSourceQuantumQueryJSON() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Terraform 0.12 supports nested maps. Ex: jsondecode(var.my_variable).attribute1.attribute2",
		ReadContext:        dataSourceQuantumQueryJSONRead, // use ReadContext for SDK v2

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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"result_map": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceQuantumQueryJSONRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	jsonStr := d.Get("json").(string)
	query := d.Get("query").(string)
	queryResult := gjson.Get(jsonStr, query)

	// Compute an ID based on the hashes of jsonStr and query.
	d.SetId(fmt.Sprintf("%d-%d", computeHash(jsonStr), computeHash(query)))
	d.Set("result", queryResult.String())

	if queryResult.IsArray() {
		var resultList []string
		for _, value := range queryResult.Array() {
			resultList = append(resultList, value.String())
		}
		if err := d.Set("result_list", resultList); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if queryResult.IsObject() {
		resultMap := make(map[string]string)
		for key, value := range queryResult.Map() {
			resultMap[key] = value.String()
		}
		if err := d.Set("result_map", resultMap); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
