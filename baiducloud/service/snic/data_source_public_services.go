package snic

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/endpoint"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func DataSourcePublicServices() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query public services that SNIC(Service Network Interface Card) can attach. \n\n",

		Read: dataSourcePublicServicesRead,

		Schema: map[string]*schema.Schema{
			"services": {
				Type:        schema.TypeList,
				Description: "Public service domain name list.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePublicServicesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithSNICClient(func(client *endpoint.Client) (interface{}, error) {
		return client.GetServices()
	})
	log.Printf("[DEBUG] Read SNIC Public Services result: %+v", raw)
	if err != nil {
		return fmt.Errorf("error reading SNIC Public Services: %w", err)
	}

	response := raw.(*endpoint.ListServiceResult)

	if err := d.Set("services", flex.FlattenStringValueList(response.Services)); err != nil {
		return fmt.Errorf("error setting services: %w", err)
	}
	d.SetId(resource.UniqueId())
	return nil
}
