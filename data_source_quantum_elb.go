package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

/*
# Usage example:
Find only ELB matching the given tags with at least one instance healthy
data "quantum_elb" "k8s_elb" {
	tags = [{ "Key" : "kubernetes.io/service-name" , "Value" : "namespace/app-elb"}]
	healthy = true
	match_all_tags = false
}
---
Find all ELB matching one of the given tags
data "quantum_elb" "k8s_elb" {
	tags = [
		{ "Key" : "kubernetes.io/service-name" , "Value" : "namespace/app-elb"},
		{ "Key" : "KubernetesCluster" , "Value" : "k8s.dev.corp"}
	]
	healthy = false
	match_all_tags = false
}
---
Find ELB matching all given tags with at least one healthy instance
data "quantum_elb" "k8s_elb" {
	tags = [
		{ "Key" : "kubernetes.io/service-name" , "Value" : "namespace/app-elb"},
		{ "Key" : "KubernetesCluster" , "Value" : "k8s.dev.corp"}
	]
	healthy = true
	match_all_tags = True
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
	var creds = credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
		})
	var awsConfig = aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(endpoints.UsEast1RegionID)
	var sess = session.New(awsConfig)
	var elbconn = elb.New(sess)
	//elbconn := meta.awsProvider.elbconn
	input := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(elbName),
	}
	describeHealth, err := elbconn.DescribeInstanceHealth(input)
	if err != nil {
		errwrap.Wrapf("Error retrieving ELB health: {{err}}", err)
	}
	for _, i := range describeHealth.InstanceStates {
		if *i.State == "InService" {
			return true
		}
	}
	return false
}

func dataSourceQuantumElb() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQuantumElbRead,
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"Key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"Value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"match_all_tags": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"healthy": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"dns_names": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceQuantumElbRead(d *schema.ResourceData, m interface{}) error {
	var creds = credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
		})
	var awsConfig = aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(endpoints.UsEast1RegionID)
	var sess = session.New(awsConfig)
	var elbconn = elb.New(sess)

	// Construct a map containing all tags
	givenTags := d.Get("tags").([]interface{})
	var elbTags []map[string]string
	for _, t := range givenTags {
		tempKey := t.(map[string]interface{})
		elbTags = append(elbTags, map[string]string{"Key": tempKey["Key"].(string), "Value": tempKey["Value"].(string)})
	}

	onlyHealthy := d.Get("healthy").(bool)
	matchAllTags := d.Get("match_all_tags").(bool)

	// Retrieve all ELB
	describeElbOpts := &elb.DescribeLoadBalancersInput{}
	describeResp, err := elbconn.DescribeLoadBalancers(describeElbOpts)
	if err != nil {
		return errwrap.Wrapf("Error retrieving ELB: {{err}}", err)
	}

	// Retrieve tags for ELB
	// In order to reduce API call we build packet of 20 ELB before asking their tags
	lbDict := make(map[string]string)
	i := 0
	var result []string
	for i < len(describeResp.LoadBalancerDescriptions) {
		lbNames := []*string{}
		// Build packet
		for _, k := range describeResp.LoadBalancerDescriptions[i : i+19] {
			lbNames = append(lbNames, k.LoadBalancerName)
			lbDict[*k.LoadBalancerName] = *k.DNSName
		}
		// Request Tags
		inputTag := &elb.DescribeTagsInput{
			LoadBalancerNames: lbNames,
		}
		tagResult, err := elbconn.DescribeTags(inputTag)
		if err != nil {
			return errwrap.Wrapf("Error retrieving ELB Tags: {{err}}", err)
		}
		i += 20

		// Check tags and ELB matching
		for _, tagDesc := range tagResult.TagDescriptions {
			if findElbTag(tagDesc, elbTags, matchAllTags) {
				if onlyHealthy {
					if isHealthy(*tagDesc.LoadBalancerName) {
						result = append(result, lbDict[*tagDesc.LoadBalancerName])
					}
				}
			}
		}
	}

	d.Set("dns_names", result)
	d.SetId("-")
	return nil
}
