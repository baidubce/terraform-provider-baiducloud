package snic

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/endpoint"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func DataSourceSNICs() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query SNIC (Service Network Interface Card) list. \n\n",

		Read: dataSourceSNICsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the vpc to which the snic belongs.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Filter by name of the snic.",
				Optional:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "Filter by ID of the subnet where the snic is located.",
				Optional:    true,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Description: "Filter by IP address of snic.",
				Optional:    true,
			},
			"service": {
				Type:        schema.TypeString,
				Description: "Filter by attached service of the snic.",
				Optional:    true,
			},
			"status": {
				Type:         schema.TypeString,
				Description:  "Filter by status of the snic. Valid valus: `available`, `unavailable`.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"available", "unavailable"}, false),
			},
			"snics": {
				Type:        schema.TypeList,
				Description: "The snic list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snic_id": {
							Type:        schema.TypeString,
							Description: "ID of the snic.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the snic.",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "ID of the vpc to which the snic belongs.",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "ID of the subnet where the snic is located.",
							Computed:    true,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Description: "IP address of the snic.",
							Computed:    true,
						},
						"service": {
							Type:        schema.TypeString,
							Description: "Attached service of the snic.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the snic.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the snic. Possible valus: `available`, `unavailable`.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSNICsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	result := []endpoint.Endpoint{}
	marker := ""

	for {
		raw, err := conn.WithSNICClient(func(client *endpoint.Client) (interface{}, error) {
			return client.ListEndpoints(buildListArgs(d, marker))
		})
		log.Printf("[DEBUG] Read SNICs result: %+v", raw)

		if err != nil {
			return fmt.Errorf("error reading SNIC list: %w", err)
		}
		response := raw.(*endpoint.ListEndpointResult)
		result = append(result, response.Endpoints...)
		if response.IsTruncated {
			marker = response.NextMarker
		} else {
			break
		}
	}

	if err := d.Set("snics", flattenSNICs(result)); err != nil {
		return fmt.Errorf("error setting snics: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func buildListArgs(d *schema.ResourceData, marker string) *endpoint.ListEndpointArgs {
	return &endpoint.ListEndpointArgs{
		VpcId:     d.Get("vpc_id").(string),
		Name:      d.Get("name").(string),
		SubnetId:  d.Get("subnet_id").(string),
		IpAddress: d.Get("ip_address").(string),
		Service:   d.Get("service").(string),
		Status:    d.Get("status").(string),
		Marker:    marker,
		MaxKeys:   1000,
	}
}

func flattenSNICs(snics []endpoint.Endpoint) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range snics {
		tfList = append(tfList, map[string]interface{}{
			"snic_id":     v.EndpointId,
			"name":        v.Name,
			"vpc_id":      v.VpcId,
			"subnet_id":   v.SubnetId,
			"ip_address":  v.IpAddress,
			"service":     v.Service,
			"description": v.Description,
			"status":      v.Status,
		})
	}
	return tfList
}
