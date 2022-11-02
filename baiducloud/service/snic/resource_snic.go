package snic

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/endpoint"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceSNIC() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage SNIC (Service Network Interface Card). \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/VPC/s/zkkus2uf2). \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceSNICCreate,
		Read:   resourceSNICRead,
		Update: resourceSNICUpdate,
		Delete: resourceSNICDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the snic. Can include letters, numbers, and `-`, `_`,`/`, `.`. Must start with a letter. The length should be between 1-65",
				Required:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the vpc to which the snic belongs.",
				Required:    true,
				ForceNew:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "ID of the subnet where the snic is located.",
				Required:    true,
				ForceNew:    true,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Description: "IP address of snic. If empty, system will automatically assign one.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"service": {
				Type: schema.TypeString,
				Description: "Attached service of the snic. Must be one of the services returned by data source `baiducloud_snic_public_services`, " +
					"and if there is descriptive text, it also needs to be included.",
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the snic.",
				Optional:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the snic. Possible valus: `available`, `unavailable`.",
				Computed:    true,
			},
		},
	}
}

func resourceSNICCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	input := buildCreationArgs(d)

	raw, err := conn.WithSNICClient(func(client *endpoint.Client) (interface{}, error) {
		return client.CreateEndpoint(input)
	})
	log.Printf("[DEBUG] Create SNIC input: %+v, result: %+v", input, raw)
	if err != nil {
		return fmt.Errorf("error creating SNIC (%s): %w", input.Name, err)
	}

	result := raw.(*endpoint.CreateEndpointResult)
	d.SetId(result.Id)

	if _, err := waitSNICAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for SNIC (%s) to become available: %w", d.Id(), err)
	}

	return resourceSNICRead(d, meta)
}

func resourceSNICRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	detail, err := FindSNIC(conn, d.Id())
	log.Printf("[DEBUG] Read SNIC (%s) result: %+v", d.Id(), detail)
	if err != nil {
		return fmt.Errorf("error reading SNIC (%s): %w", d.Id(), err)
	}

	if err := d.Set("name", detail.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("vpc_id", detail.VpcId); err != nil {
		return fmt.Errorf("error setting vpc_id: %w", err)
	}
	if err := d.Set("subnet_id", detail.SubnetId); err != nil {
		return fmt.Errorf("error setting subnet_id: %w", err)
	}
	if err := d.Set("ip_address", detail.IpAddress); err != nil {
		return fmt.Errorf("error setting ip_address: %w", err)
	}
	if err := d.Set("service", detail.Service); err != nil {
		return fmt.Errorf("error setting service: %w", err)
	}
	if err := d.Set("description", detail.Description); err != nil {
		return fmt.Errorf("error setting description: %w", err)
	}
	if err := d.Set("status", detail.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	return nil
}

func resourceSNICUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	_, err := conn.WithSNICClient(func(client *endpoint.Client) (interface{}, error) {
		args := &endpoint.UpdateEndpointArgs{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		return nil, client.UpdateEndpoint(d.Id(), args)
	})
	log.Printf("[DEBUG] Update SNIC (%s)", d.Id())
	if err != nil {
		return fmt.Errorf("error updating SNIC (%s): %w", d.Id(), err)
	}
	return resourceSNICRead(d, meta)
}

func resourceSNICDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	_, err := conn.WithSNICClient(func(client *endpoint.Client) (interface{}, error) {
		return nil, client.DeleteEndpoint(d.Id(), "")
	})
	log.Printf("[DEBUG] Delete SNIC (%s)", d.Id())
	if err != nil {
		return fmt.Errorf("error deleting SNIC (%s): %w", d.Id(), err)
	}

	return nil
}

func buildCreationArgs(d *schema.ResourceData) *endpoint.CreateEndpointArgs {
	input := &endpoint.CreateEndpointArgs{}

	input.Name = d.Get("name").(string)
	input.VpcId = d.Get("vpc_id").(string)
	input.SubnetId = d.Get("subnet_id").(string)
	input.IpAddress = d.Get("ip_address").(string)
	input.Service = d.Get("service").(string)
	input.Description = d.Get("description").(string)
	input.Billing = &endpoint.Billing{
		PaymentTiming: "Postpaid",
	}

	return input
}
