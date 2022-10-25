package cdn

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func ResourceDomainConfigOrigin() *schema.Resource {

	return &schema.Resource{

		Description: "Use this resource to manage Forward-to-origin configuration of the acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN/s/xjxzi7729). \n\n" +
			"~> **NOTE:** Creating a resource will overwrite current Forward-to-origin configuration. " +
			"Deleting a resource won't change current configuration.",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceDomainConfigOriginCreate,
		Read:   resourceDomainConfigOriginRead,
		Update: resourceDomainConfigOriginUpdate,
		Delete: resourceDomainConfigOriginDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Name of the acceleration domain.",
				Required:    true,
				ForceNew:    true,
			},
			"range_switch": {
				Type:         schema.TypeString,
				Description:  "Whether range forwarding to origin is enabled. Defaults to `off`. Other valid value: `on`",
				Optional:     true,
				Default:      "off",
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"origin_protocol": {
				Type:        schema.TypeList,
				Description: "Forward-to-origin protocol configuration of the acceleration domain.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeString,
							Description:  "Protocol value. Defaults to `http`. Other valid values: `https`, `*`",
							Optional:     true,
							Default:      "http",
							ValidateFunc: validation.StringInSlice([]string{"*", "http", "https"}, false),
						},
					},
				},
			},
			"offline_mode": {
				Type:        schema.TypeBool,
				Description: "Whether offline mode is enabled. Defaults to `false`.",
				Optional:    true,
				Default:     false,
			},
			"client_ip": {
				Type:        schema.TypeList,
				Description: "Getting user's real IP configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether getting user's real IP is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"name": {
							Type:         schema.TypeString,
							Description:  "IP type. Valid values: `True-Client-Ip`, `X-Real-IP`.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"True-Client-Ip", "X-Real-IP"}, false),
						},
					},
				},
				DiffSuppressFunc: clientIpDiffSuppress,
			},
		},
	}
}

func resourceDomainConfigOriginCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Get("domain").(string)

	if err := updateConfigOrigin(d, conn, domain); err != nil {
		return err
	}

	d.SetId(domain)
	return resourceDomainConfigOriginRead(d, meta)
}

func resourceDomainConfigOriginRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()
	d.Set("domain", domain)

	if err := readCommonConfigOrigin(d, conn, domain); err != nil {
		return err
	}
	if err := readOriginProtocol(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func resourceDomainConfigOriginUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	if err := updateConfigOrigin(d, conn, domain); err != nil {
		return err
	}
	return resourceDomainConfigOriginRead(d, meta)
}

func resourceDomainConfigOriginDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func readCommonConfigOrigin(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	config, err := FindDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config: %w", domain, err)
	}
	log.Printf("[DEBUG] Read CDN Domain (%s) Config RangeSwitch result: %+v", domain, config.RangeSwitch)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config OfflineMode result: %+v", domain, config.OfflineMode)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config ClientIp result: %+v", domain, config.ClientIp)
	//
	d.Set("range_switch", config.RangeSwitch)
	d.Set("offline_mode", config.OfflineMode)
	d.Set("client_ip", flattenClientIp(config.ClientIp))
	return nil
}

func readOriginProtocol(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetOriginProtocol(domain)
	})
	log.Printf("[DEBUG] Read CDN Domain (%s) Config OriginProtocol result: %+v", domain, raw)

	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config OriginProtocol: %w", domain, err)
	}

	d.Set("origin_protocol", flattenOriginProtocol(raw.(string)))
	return nil
}

func updateConfigOrigin(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateRangeSwitch(d, conn, domain); err != nil {
		return err
	}
	if err := updateOriginProtocol(d, conn, domain); err != nil {
		return err
	}
	if err := updateOfflineMode(d, conn, domain); err != nil {
		return err
	}
	if err := updateClientIp(d, conn, domain); err != nil {
		return err
	}

	return nil
}

func updateRangeSwitch(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("range_switch") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config RangeSwitch", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetRangeSwitch(domain, d.Get("range_switch").(string) == "on")
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config RangeSwitch: %w", domain, err)
		}
	}
	return nil
}

func updateOriginProtocol(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("origin_protocol") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config OriginProtocol", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetOriginProtocol(domain, expandOriginProtocol(d.Get("origin_protocol").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config OriginProtocol: %w", domain, err)
		}
	}
	return nil
}

func updateOfflineMode(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("offline_mode") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config OfflineMode", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetOfflineMode(domain, d.Get("offline_mode").(bool))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config OfflineMode: %w", domain, err)
		}
	}
	return nil
}

func updateClientIp(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("client_ip") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config ClientIp", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetClientIp(domain, expandClientIp(d.Get("client_ip").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config ClientIp: %w", domain, err)
		}
	}
	return nil
}
