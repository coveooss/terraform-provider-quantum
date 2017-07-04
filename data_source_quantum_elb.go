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
	session := getELBSession()

	// Construct a list containing all tags
	var elbTags []keyValue
	for _, t := range d.Get("tags").([]interface{}) {
		tempKey := t.(map[string]interface{})
		elbTags = append(elbTags, keyValue{tempKey["Key"].(string), tempKey["Value"].(string)})
	}

	onlyHealthy := d.Get("healthy").(bool)
	matchAllTags := d.Get("match_all_tags").(bool)

	// Retrieve all ELB
	describeResp, err := session.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})
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
		endRange := 0
		if i+19 > len(describeResp.LoadBalancerDescriptions) {
			endRange = len(describeResp.LoadBalancerDescriptions)
		} else {
			endRange = i + 19
		}
		for _, k := range describeResp.LoadBalancerDescriptions[i:endRange] {
			lbNames = append(lbNames, k.LoadBalancerName)
			lbDict[*k.LoadBalancerName] = *k.DNSName
		}
		// Request Tags
		tagResult, err := session.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: lbNames})
		if err != nil {
			return errwrap.Wrapf("Error retrieving ELB Tags: {{err}}", err)
		}
		i += 20

		// Check tags and ELB matching
		for _, tagDesc := range tagResult.TagDescriptions {
			if findElbTag(tagDesc, elbTags, matchAllTags) {
				if onlyHealthy {
					if isHealthy(session, *tagDesc.LoadBalancerName) {
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

func getELBSession() *elb.ELB {
	var creds = credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
		})
	var awsConfig = aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(endpoints.UsEast1RegionID)
	return elb.New(session.New(awsConfig))
}

func isHealthy(session *elb.ELB, elbName string) bool {
	input := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(elbName),
	}
	describeHealth, err := session.DescribeInstanceHealth(input)
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

func findElbTag(elbTags *elb.TagDescription, queryTags []keyValue, matchAll bool) bool {
	matchCount := 0
	for _, key := range elbTags.Tags {
		for _, tag := range queryTags {
			if *key.Key == tag.key && *key.Value == tag.value {
				if !matchAll {
					return true
				}
				matchCount++
			}
		}
	}
	return matchCount == len(queryTags)
}

type keyValue struct {
	key   string
	value string
}
