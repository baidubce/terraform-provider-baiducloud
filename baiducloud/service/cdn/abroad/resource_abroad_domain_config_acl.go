package abroad

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/cdn/abroad"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"

	"log"
)

func ResourceAbroadDomainConfigACL() *schema.Resource {

	return &schema.Resource{

		Description: "Use this resource to manage acl-related configuration of the abroad acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN-ABROAD/s/ekbsxow69). \n\n" +
			"~> **NOTE:** Creating a resource will overwrite current acl-related configuration. " +
			"Deleting a resource won't change current configuration.",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceDomainConfigACLCreate,
		Read:   resourceDomainConfigACLRead,
		Update: resourceDomainConfigACLUpdate,
		Delete: resourceDomainConfigACLDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Name of the acceleration domain.",
				Required:    true,
				ForceNew:    true,
			},
			"allow_empty": {
				Type:        schema.TypeBool,
				Description: "Whether empty referer access is allowed. Defaults to `true`.",
				Optional:    true,
				Default:     true,
			},
			"referer_acl": {
				Type:        schema.TypeList,
				Description: "Referer access configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"black_list": {
							Type:        schema.TypeSet,
							Description: "Referer blacklist. Support wildcard and no protocol required. Conflict with `white_list`",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ConflictsWith: []string{"referer_acl.0.white_list"},
						},
						"white_list": {
							Type:        schema.TypeSet,
							Description: "Referer whitelist. Support wildcard and no protocol required. Conflict with `black_list`",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ConflictsWith: []string{"referer_acl.0.black_list"},
						},
					},
				},
			},
			"ip_acl": {
				Type:        schema.TypeList,
				Description: "IP access configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"black_list": {
							Type:        schema.TypeSet,
							Description: "IP blacklist, support IP segments in CIDR format. Conflict with `white_list`",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ConflictsWith: []string{"ip_acl.0.white_list"},
						},
						"white_list": {
							Type:        schema.TypeSet,
							Description: "IP whitelist, support IP segments in CIDR format. Conflict with `black_list`",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ConflictsWith: []string{"ip_acl.0.black_list"},
						},
					},
				},
			},
		},
	}
}

func resourceDomainConfigACLCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Get("domain").(string)

	if err := updateConfigACL(d, conn, domain); err != nil {
		return err
	}

	d.SetId(domain)
	// wait for running status
	if _, err := waitAbroadCDNDomainAvailable(conn, domain); err != nil {
		return fmt.Errorf("error waiting Abraod CDN domain (%s) becoming available: %w", d.Id(), err)
	}
	return resourceDomainConfigACLRead(d, meta)
}

func resourceDomainConfigACLRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()
	err := d.Set("domain", domain)
	if err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) Config domain: %w", domain, err)
	}

	if err := readCommonConfigACL(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func resourceDomainConfigACLUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	if err := updateConfigACL(d, conn, domain); err != nil {
		return err
	}
	// wait for running status
	if _, err := waitAbroadCDNDomainAvailable(conn, domain); err != nil {
		return fmt.Errorf("error waiting Abraod CDN domain (%s) becoming available: %w", d.Id(), err)
	}
	return resourceDomainConfigACLRead(d, meta)
}

func resourceDomainConfigACLDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func readCommonConfigACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	config, err := FindAbroadDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error getting abroad CDN Domain (%s) Config: %w", domain, err)
	}
	log.Printf("[DEBUG] Read abroad CDN Domain (%s) Config RefererACL result: %+v", domain, config.RefererACL)
	log.Printf("[DEBUG] Read abroad CDN Domain (%s) Config IpACL result: %+v", domain, config.IpACL)

	err = d.Set("referer_acl", flattenRefererACL(config.RefererACL))
	if err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) Config RefererACL: %w", domain, err)
	}
	err = d.Set("ip_acl", flattenIpACL(config.IpACL))
	if err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) Config ip_acl: %w", domain, err)
	}
	err = d.Set("allow_empty", config.RefererACL.AllowEmpty)
	if err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) Config allow_empty: %w", domain, err)
	}

	return nil
}

func updateConfigACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateRefererACL(d, conn, domain); err != nil {
		return err
	}
	if err := updateIpACL(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func updateRefererACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.HasChange("referer_acl") {
		expandRefererACL := expandRefererACL(d.Get("referer_acl").([]interface{}))
		log.Printf("[DEBUG] Update abroad CDN Domain (%s) Config RefererACL: %+v", domain, expandRefererACL)
		expandRefererACL.AllowEmpty = d.Get("allow_empty").(bool)
		_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
			return nil, client.SetRefererACL(domain, expandRefererACL)
		})
		if err != nil {
			return fmt.Errorf("error updating abroad CDN Domain (%s) Config RefererACL: %w", domain, err)
		}
	}
	return nil
}

func updateIpACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.HasChange("ip_acl") {
		ipACL := expandIpACL(d.Get("ip_acl").([]interface{}))
		log.Printf("[DEBUG] Update abroad CDN Domain (%s) Config IP Acl: %+v", domain, ipACL)
		_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
			return nil, client.SetIpACL(domain, ipACL)
		})
		if err != nil {
			return fmt.Errorf("error updating abroad CDN Domain (%s) Config Ip ACL: %w", domain, err)
		}
	}
	return nil
}
