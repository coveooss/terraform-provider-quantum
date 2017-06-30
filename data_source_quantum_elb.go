package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

/*
# Usage example:
Find only ELB matching the given tags with at least one instance healthy
data "quantum_elb" "k8s_elb" {
	tag = [{ "Key" : "kubernetes.io/service-name" , "Value" : "namespace/app-elb"}]
	healthy = true
	matchAllTags = false
}
---
Find all ELB matching one of the given tags 
data "quantum_elb" "k8s_elb" {
	tag = [
		{ "Key" : "kubernetes.io/service-name" , "Value" : "namespace/app-elb"},
		{ "Key" : "KubernetesCluster" , "Value" : "k8s.dev.corp"}
	]
	healthy = false
	matchAllTags = false
}
---
Find ELB matching all given tags with at least one healthy instance
data "quantum_elb" "k8s_elb" {
	tag = [
		{ "Key" : "kubernetes.io/service-name" , "Value" : "namespace/app-elb"},
		{ "Key" : "KubernetesCluster" , "Value" : "k8s.dev.corp"}
	]
	healthy = true
	matchAllTags = True
}
*/

func findElbTag(elbTags *elb.TagDescription, queryTags []map[string]string, matchAll bool) bool {
	matchCount := 0
	for _, key := range elbTags.Tags {
		for _, tag := range queryTags {
			if *key.Key == tag["Key"] && *key.Value == tag["Value"] {
				if matchAll {
					matchCount++
					if matchCount == len(queryTags) {
						return true
					}
				} else {
					return true
				}
			}
		}
	}
	return false
}

func isHealthy(elbName string) bool {
	elbconn := meta.(*AWSClient).elbconn
	input := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(elbName),
	}
	describeHealth, err := elbconn.DescribeInstanceHealth(input)
	if err != nil {
		errwrap.Wrapf("Error retrieving ELB health: {{err}}", err)
	}
	for _, i := range describeHealth.InstanceStates {
		if *i.State == "Healthy" {
			return true
		}
	}
	return false
}

func dataSourceQuantumElb() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQuantumElbRead,
		Schema: map[string]*schema.Schema{
			"tag": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"Key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"Value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"matchAllTags": {
				Type: schema.TypeBool,
			},
			"healthy": {
				Type: schema.TypeBool,
			},
		},
	}
}

func dataSourceQuantumElbRead(d *schema.ResourceData, m interface{}) error {
	elbconn := meta.(*AWSClient).elbconn
	elbTag := d.Get("tag").(map[string]string)
	onlyHealthy := d.Get("healthy").(bool)
	matchAllTags := d.Get(":= d.Get("healthy").(bool)").(bool)

	// Retrieve all ELB
	describeElbOpts := &elb.DescribeLoadBalancersInput{}
	describeResp, err := elbconn.DescribeLoadBalancers(describeElbOpts)
	if err != nil {
		return errwrap.Wrapf("Error retrieving ELB: {{err}}", err)
	}

	// Retrieve tags for ELB
	// In order to reduce API call we build packet of 20 ELB before asking their tags
	lbDict := make(map[string]elb.LoadBalancerDescription)
	i := 0
	var result []string
	for i < len(describeResp.LoadBalancerDescriptions) {
		lbNames := []*string{}
		for _, k := range describeResp.LoadBalancerDescriptions[i : i+19] {
			lbNames = append(lbNames, k.LoadBalancerName)
			lbDict[*k.LoadBalancerName] = *k.DNSName
		}
		inputTag := &elb.DescribeTagsInput{
			LoadBalancerNames: lbNames,
		}
		tagResult, err := elbconn.DescribeTags(inputTag)
		if err != nil {
			return errwrap.Wrapf("Error retrieving ELB Tags: {{err}}", err)
		}
		i += 20

		for _, tagDesc := range tagResult.TagDescriptions {
			if findElbTag(tagDesc, elbTag, matchAllTags) {
				if onlyHealthy {
					if isHealthy(*tagDesc.LoadBalancerName) {
						result = append(result, lbDict[*tagDesc.LoadBalancerName])
					}
				}
			}
		}
	}

	d.Set("dnsNames", result)
	d.SetId("-")
	return nil
}
