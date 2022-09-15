package cdn

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"reflect"
)

func ResourceDomainConfigACL() *schema.Resource {

	return &schema.Resource{

		Description: "Use this resource to manage acl-related configuration of the acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN/s/yjxzhvf21). \n\n" +
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
						"allow_empty": {
							Type:        schema.TypeBool,
							Description: "Whether empty referer access is allowed. Defaults to `true`.",
							Optional:    true,
							Default:     true,
						},
					},
				},
				DiffSuppressFunc: refererACLDiffSuppress,
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
				DiffSuppressFunc: ipACLDiffSuppress,
			},
			"ua_acl": {
				Type:        schema.TypeList,
				Description: "User agent access configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"black_list": {
							Type:        schema.TypeSet,
							Description: "UA blacklist, length of a single ua should be in 1-200.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ConflictsWith: []string{"ua_acl.0.white_list"},
						},
						"white_list": {
							Type:        schema.TypeSet,
							Description: "UA whitelist, length of a single ua should be in 1-200.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ConflictsWith: []string{"ua_acl.0.black_list"},
						},
					},
				},
				DiffSuppressFunc: uaACLDiffSuppress,
			},
			"cors": {
				Type:        schema.TypeList,
				Description: "Cors cross-domain configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:         schema.TypeString,
							Description:  "Whether allowing cross-domain access. Defaults to `off`. Other valid value: `on`",
							Optional:     true,
							Default:      "off",
							ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
						},
						"origin_list": {
							Type:        schema.TypeSet,
							Description: "Domain names with cross-domain allowed. Support extensive domain name. Each name can contain at most one wildcard",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				DiffSuppressFunc: corsDiffSuppress,
			},
			"access_limit": {
				Type:        schema.TypeList,
				Description: "IP access limit configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether visit frequency limit for a single IP node is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"limit": {
							Type:         schema.TypeInt,
							Description:  "Maximum number of requests a single IP node can send in one second. Defaults to `1000`.",
							Optional:     true,
							Default:      1000,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
				DiffSuppressFunc: accessLimitDiffSuppress,
			},
			"traffic_limit": {
				Type:        schema.TypeList,
				Description: "Rate limit for a single link configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Description: "Whether rate limit is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"limit_rate": {
							Type:        schema.TypeInt,
							Description: "Limit rate in Byte/s",
							Optional:    true,
						},
						"limit_start_hour": {
							Type:         schema.TypeInt,
							Description:  "Time to start speed limit. Should be in 0-24, and smaller than `limit_end_hour`. Defaults to `0`.",
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 24),
						},
						"limit_end_hour": {
							Type:         schema.TypeInt,
							Description:  "Time to start speed limit. Should be in 0-24, and greater than `limit_start_hour`. Defaults to `24`.",
							Optional:     true,
							Default:      24,
							ValidateFunc: validation.IntBetween(0, 24),
						},
					},
				},
				DiffSuppressFunc: trafficLimitDiffSuppress,
			},
			"request_auth": {
				Type:        schema.TypeList,
				Description: "Access authentication configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Description:  "Authentication method. Valid values: `A`, `B`, `C`",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"A", "B", "C"}, false),
						},
						"key1": {
							Type:        schema.TypeString,
							Description: "Main authorization key. Letters and numbers can be used. Length should be in 6-32.",
							Required:    true,
						},
						"key2": {
							Type:        schema.TypeString,
							Description: "Secondary authorization key. Letters and numbers can be used. Length should be in 6-32.",
							Optional:    true,
						},
						"timeout": {
							Type:        schema.TypeInt,
							Description: "Authorization cache time.",
							Optional:    true,
						},
						"timestamp_metric": {
							Type:         schema.TypeInt,
							Description:  "Time format. Valid values: `10`, `16`.",
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{10, 16}),
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
	return nil
}

func resourceDomainConfigACLRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	d.Set("domain", domain)

	if err := readRefererACL(d, conn, domain); err != nil {
		return err
	}
	if err := readIpACL(d, conn, domain); err != nil {
		return err
	}
	if err := readUaACL(d, conn, domain); err != nil {
		return err
	}
	if err := readCors(d, conn, domain); err != nil {
		return err
	}
	if err := readAccessLimit(d, conn, domain); err != nil {
		return err
	}
	if err := readTrafficLimit(d, conn, domain); err != nil {
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
	return resourceDomainConfigACLRead(d, meta)
}

func resourceDomainConfigACLDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func updateConfigACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateRefererACL(d, conn, domain); err != nil {
		return err
	}
	if err := updateIpACL(d, conn, domain); err != nil {
		return err
	}
	if err := updateUaACL(d, conn, domain); err != nil {
		return err
	}
	if err := updateCors(d, conn, domain); err != nil {
		return err
	}
	if err := updateAccessLimit(d, conn, domain); err != nil {
		return err
	}
	if err := updateTrafficLimit(d, conn, domain); err != nil {
		return err
	}
	if err := updateRequestAuth(d, conn, domain); err != nil {
		return err
	}

	return nil
}

//<editor-fold desc="RefererACL">
func readRefererACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	log.Printf("[DEBUG] Read CDN Domain Config RefererACL: %s", domain)

	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetRefererACL(domain)
	})
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config RefererACL: %w", domain, err)
	}

	d.Set("referer_acl", flattenRefererACL(raw.(*api.RefererACL)))
	return nil
}

func updateRefererACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("referer_acl") {
		oldV, newV := d.GetChange("referer_acl")
		oldRefererACL := expandRefererACL(oldV.([]interface{}))
		newRefererACL := expandRefererACL(newV.([]interface{}))
		if reflect.DeepEqual(oldRefererACL, newRefererACL) {
			return nil
		}

		log.Printf("[DEBUG] Update CDN Domain Config RefererACL: %s", domain)

		allowEmpty := newRefererACL.AllowEmpty
		allowEmptyUpdated := false
		needUpdateAllowEmpty := oldRefererACL.AllowEmpty != newRefererACL.AllowEmpty

		if len(oldRefererACL.BlackList) > 0 && len(newRefererACL.BlackList) == 0 {
			if err := setRefererACL(conn, domain, []string{}, true, allowEmpty); err != nil {
				return err
			}
			allowEmptyUpdated = true
		}
		if len(oldRefererACL.WhiteList) > 0 && len(newRefererACL.WhiteList) == 0 {
			if err := setRefererACL(conn, domain, []string{}, false, allowEmpty); err != nil {
				return err
			}
			allowEmptyUpdated = true
		}
		if len(newRefererACL.BlackList) > 0 {
			if err := setRefererACL(conn, domain, newRefererACL.BlackList, true, allowEmpty); err != nil {
				return err
			}
			allowEmptyUpdated = true
		}
		if len(newRefererACL.WhiteList) > 0 || (needUpdateAllowEmpty && !allowEmptyUpdated) {
			if err := setRefererACL(conn, domain, newRefererACL.WhiteList, false, allowEmpty); err != nil {
				return err
			}
		}

	}
	return nil
}

func setRefererACL(conn *connectivity.BaiduClient, domain string, list []string, setBlackList bool, allowEmpty bool) error {
	_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		if setBlackList {
			return nil, client.SetRefererACL(domain, list, nil, allowEmpty)
		} else {
			return nil, client.SetRefererACL(domain, nil, list, allowEmpty)
		}
	})
	if err != nil {
		return fmt.Errorf("error updating CDN Domain (%s) Config RefererACL: %w", domain, err)
	}
	return nil
}

//</editor-fold>

//<editor-fold desc="IpACL">
func readIpACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	log.Printf("[DEBUG] Read CDN Domain Config IpACL: %s", domain)

	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetIpACL(domain)
	})
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config IpACL: %w", domain, err)
	}

	d.Set("ip_acl", flattenIpACL(raw.(*api.IpACL)))
	return nil
}

func updateIpACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("ip_acl") {
		oldV, newV := d.GetChange("ip_acl")
		oldIpACL := expandIpACL(oldV.([]interface{}))
		newIpACL := expandIpACL(newV.([]interface{}))
		if reflect.DeepEqual(oldIpACL, newIpACL) {
			return nil
		}

		log.Printf("[DEBUG] Update CDN Domain Config IpACL: %s", domain)

		if len(oldIpACL.BlackList) > 0 && len(newIpACL.BlackList) == 0 {
			if err := setIpACL(conn, domain, []string{}, true); err != nil {
				return err
			}
		}
		if len(oldIpACL.WhiteList) > 0 && len(newIpACL.WhiteList) == 0 {
			if err := setIpACL(conn, domain, []string{}, false); err != nil {
				return err
			}
		}
		if len(newIpACL.BlackList) > 0 {
			if err := setIpACL(conn, domain, newIpACL.BlackList, true); err != nil {
				return err
			}
		}
		if len(newIpACL.WhiteList) > 0 {
			if err := setIpACL(conn, domain, newIpACL.WhiteList, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func setIpACL(conn *connectivity.BaiduClient, domain string, list []string, setBlackList bool) error {
	_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		if setBlackList {
			return nil, client.SetIpACL(domain, list, nil)
		} else {
			return nil, client.SetIpACL(domain, nil, list)
		}
	})
	if err != nil {
		return fmt.Errorf("error updating CDN Domain (%s) Config IpACL: %w", domain, err)
	}
	return nil
}

//</editor-fold>

//<editor-fold desc="UaACL">
func readUaACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	log.Printf("[DEBUG] Read CDN Domain Config UaACL: %s", domain)

	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetUaACL(domain)
	})
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config UaACL: %w", domain, err)
	}

	d.Set("ua_acl", flattenUaACL(raw.(*api.UaACL)))
	return nil
}

func updateUaACL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("ua_acl") {
		oldV, newV := d.GetChange("ua_acl")
		oldUaACL := expandUaACL(oldV.([]interface{}))
		newUaACL := expandUaACL(newV.([]interface{}))
		if reflect.DeepEqual(oldUaACL, newUaACL) {
			return nil
		}

		log.Printf("[DEBUG] Update CDN Domain Config UaACL: %s", domain)

		if len(oldUaACL.BlackList) > 0 && len(newUaACL.BlackList) == 0 {
			if err := setUaACL(conn, domain, []string{}, true); err != nil {
				return err
			}
		}
		if len(oldUaACL.WhiteList) > 0 && len(newUaACL.WhiteList) == 0 {
			if err := setUaACL(conn, domain, []string{}, false); err != nil {
				return err
			}
		}
		if len(newUaACL.BlackList) > 0 {
			if err := setUaACL(conn, domain, newUaACL.BlackList, true); err != nil {
				return err
			}
		}
		if len(newUaACL.WhiteList) > 0 {
			if err := setUaACL(conn, domain, newUaACL.WhiteList, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func setUaACL(conn *connectivity.BaiduClient, domain string, list []string, setBlackList bool) error {
	_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		if setBlackList {
			return nil, client.SetUaACL(domain, list, nil)
		} else {
			return nil, client.SetUaACL(domain, nil, list)
		}
	})
	if err != nil {
		return fmt.Errorf("error updating CDN Domain (%s) Config UaACL: %w", domain, err)
	}
	return nil
}

//</editor-fold>

//<editor-fold desc="Cors">
func readCors(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	log.Printf("[DEBUG] Read CDN Domain Config Cors: %s", domain)

	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetCors(domain)
	})
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config Cors: %w", domain, err)
	}

	d.Set("cors", flattenCors(raw.(*api.CorsCfg)))
	return nil
}

func updateCors(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("cors") {
		oldV, newV := d.GetChange("cors")
		oldCors := expandCors(oldV.([]interface{}))
		newCors := expandCors(newV.([]interface{}))
		if reflect.DeepEqual(oldCors, newCors) {
			return nil
		}

		log.Printf("[DEBUG] Update CDN Domain Config Cors: %s", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetCors(domain, newCors.IsAllow, newCors.Origins)
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config Cors: %w", domain, err)
		}
	}
	return nil
}

//</editor-fold>

//<editor-fold desc="AccessLimit">
func readAccessLimit(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	log.Printf("[DEBUG] Read CDN Domain Config AccessLimit: %s", domain)

	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetAccessLimit(domain)
	})
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config AccessLimit: %w", domain, err)
	}

	d.Set("access_limit", flattenAccessLimit(raw.(*api.AccessLimit)))
	return nil
}

func updateAccessLimit(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("access_limit") {
		log.Printf("[DEBUG] Update CDN Domain Config AccessLimit: %s", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetAccessLimit(domain, expandAccessLimit(d.Get("access_limit").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config AccessLimit: %w", domain, err)
		}
	}
	return nil
}

//</editor-fold>

//<editor-fold desc="TrafficLimit">
func readTrafficLimit(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	log.Printf("[DEBUG] Read CDN Domain Config TrafficLimit: %s", domain)

	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetTrafficLimit(domain)
	})
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config TrafficLimit: %w", domain, err)
	}

	d.Set("traffic_limit", flattenTrafficLimit(raw.(*api.TrafficLimit)))
	return nil
}

func updateTrafficLimit(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("traffic_limit") {
		log.Printf("[DEBUG] Update CDN Domain Config TrafficLimit: %s", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetTrafficLimit(domain, expandTrafficLimit(d.Get("traffic_limit").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config TrafficLimit: %w", domain, err)
		}
	}
	return nil
}

//</editor-fold>

//<editor-fold desc="RequestAuth">
func updateRequestAuth(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("request_auth") {
		log.Printf("[DEBUG] Update CDN Domain Config RequestAuth: %s", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetDomainRequestAuth(domain, expandRequestAuth(d.Get("request_auth").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config RequestAuth: %w", domain, err)
		}
	}
	return nil
}

//</editor-fold>
