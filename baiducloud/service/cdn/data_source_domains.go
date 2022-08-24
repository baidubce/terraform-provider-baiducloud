package cdn

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func DataSourceDomains() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query all domain names under the user. \n\n",

		Read: dataSourceDomainsRead,
		Schema: map[string]*schema.Schema{
			"domains": {
				Type:        schema.TypeList,
				Description: "Domain name list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Description: "Name of the acceleration domain.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the domain name. Possible values: `RUNNING`,`OPERATING`, `STOPPED`.",
							Computed:    true,
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Description:  "Domain status filter. Defaults to `ALL`. Other valid values: `RUNNING`,`OPERATING`, `STOPPED`",
				Optional:     true,
				Default:      DomainStatusAll,
				ValidateFunc: validation.StringInSlice([]string{DomainStatusAll, DomainStatusRunning, DomainStatusOperating, DomainStatusStopped}, false),
			},
			"rule": {
				Type:        schema.TypeString,
				Description: "Domain name filter. Support fuzzy matching. Can only contain letters, numbers and periods",
				Optional:    true,
			},
		},
	}
}

func dataSourceDomainsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	status := d.Get("status").(string)
	rule := d.Get("rule").(string)

	log.Printf("[DEBUG] Read CDN Domain List: status(%s), rule(%s)", status, rule)

	result, err := FindDomainsStatus(conn, status, rule)
	if err != nil {
		return fmt.Errorf("error getting CDN Domain list: %w", err)
	}

	d.SetId(resource.UniqueId())
	d.Set("domains", flattenDomainStatuses(result))

	return nil
}
