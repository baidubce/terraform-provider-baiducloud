package abroad

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/cdn/abroad"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"

	"log"
)

func ResourceAbroadDomainConfigCache() *schema.Resource {

	return &schema.Resource{

		Description: "Use this resource to manage cache-related configuration of the abroad acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN-ABROAD/s/Zkbstm0vg). \n\n" +
			"~> **NOTE:** Creating a resource will overwrite current cache-related configuration. " +
			"Deleting a resource won't change current configuration.",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceDomainConfigCacheCreate,
		Read:   resourceDomainConfigCacheRead,
		Update: resourceDomainConfigCacheUpdate,
		Delete: resourceDomainConfigCacheDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Name of the abroad acceleration domain.",
				Required:    true,
				ForceNew:    true,
			},
			"cache_ttl": {
				Type:        schema.TypeSet,
				Description: "Cache expiration rules of the abroad acceleration domain.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type: schema.TypeString,
							Description: "Cache rule type. Valid values: `suffix`(file name suffix), " +
								"`path`(directory in the url), `origin`(origin server rule. There is only one such rule," +
								" and only `weight` is required. Set `value` to `-`, `ttl` to `0`), " +
								"`exactPath`(path is completely matched).",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"suffix", "path", "exactPath"}, false),
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Configuration rule for the specified type.",
							Required:    true,
						},
						"weight": {
							Type: schema.TypeInt,
							Description: "The origin server weight. Must be between `0` and `100`. Defaults to `0`. " +
								"The higher the weight, the higher the priority. No effect when `type` is `code`.",
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"ttl": {
							Type:        schema.TypeInt,
							Description: "Cache duration in seconds.",
							Required:    true,
						},
						"override_origin": {
							Type:        schema.TypeBool,
							Description: "Whether to override the origin site’s caching rules. Defaults to `true`",
							Optional:    true,
							Default:     true,
						},
					},
				},
			},
			"cache_full_url": {
				Type:        schema.TypeBool,
				Description: "Whether caching of full url is enabled. Defaults to `true`",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceDomainConfigCacheCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Get("domain").(string)

	if err := updateConfigCache(d, conn, domain); err != nil {
		return err
	}

	d.SetId(domain)
	// wait for running status
	if _, err := waitAbroadCDNDomainAvailable(conn, domain); err != nil {
		return fmt.Errorf("error waiting Abraod CDN domain (%s) becoming available: %w", d.Id(), err)
	}
	return resourceDomainConfigCacheRead(d, meta)
}

func resourceDomainConfigCacheRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}

	if err := readCommonConfigCache(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func resourceDomainConfigCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	if err := updateConfigCache(d, conn, domain); err != nil {
		return err
	}
	// wait for running status
	if _, err := waitAbroadCDNDomainAvailable(conn, domain); err != nil {
		return fmt.Errorf("error waiting Abraod CDN domain (%s) becoming available: %w", d.Id(), err)
	}
	return resourceDomainConfigCacheRead(d, meta)
}

func resourceDomainConfigCacheDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func readCommonConfigCache(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	config, err := FindAbroadDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error getting abroad CDN Domain (%s) Config: %w", domain, err)
	}

	err = d.Set("cache_ttl", flattenAbroadCacheTTLs(config.CacheTTL))
	if err != nil {
		return fmt.Errorf("error setting abroad CDN Domain (%s) Config: %w", domain, err)
	}
	err = d.Set("cache_full_url", config.CacheFullUrl)
	if err != nil {
		return fmt.Errorf("error setting abroad CDN Domain (%s) Config: %w", domain, err)
	}

	return nil
}

func updateConfigCache(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateCacheTTL(d, conn, domain); err != nil {
		return err
	}
	if err := updateCacheFullURL(d, conn, domain); err != nil {
		return err
	}

	return nil
}

func updateCacheTTL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("cache_ttl") {
		log.Printf("[DEBUG] Update Abroad CDN Domain (%s) Config CacheTTL", domain)

		_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
			schema := d.Get("cache_ttl").(*schema.Set)
			return nil, client.SetCacheTTL(domain, expandAbroadCacheTTLs(schema))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config CacheTTL: %w", domain, err)
		}
	}
	return nil
}

func updateCacheFullURL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("cache_full_url") {
		log.Printf("[DEBUG] Update Abroad CDN Domain (%s) Config cache full url", domain)

		_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
			return nil, client.SetCacheFullUrl(domain, d.Get("cache_full_url").(bool))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config cache full url: %w", domain, err)
		}
	}
	return nil
}
