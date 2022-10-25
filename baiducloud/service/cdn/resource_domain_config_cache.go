package cdn

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func ResourceDomainConfigCache() *schema.Resource {

	return &schema.Resource{

		Description: "Use this resource to manage cache-related configuration of the acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN/s/wjxzhgxnx). \n\n" +
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
				Description: "Name of the acceleration domain.",
				Required:    true,
				ForceNew:    true,
			},
			"cache_ttl": {
				Type:        schema.TypeSet,
				Description: "Cache expiration rules of the acceleration domain.",
				Optional:    true,
				Set:         cacheTTLHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "Cache rule type. Valid values: `suffix`(file name suffix), `path`(directory in the url), `origin`(origin server rule. There is only one such rule, and only `weight` is required. Set `value` to `-`, `ttl` to `0`), `code`(status code cache, currently supports `404`, `502`, `503`, `504`), `exactPath`(path is completely matched).",
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Configuration rule for the specified type.",
							Required:    true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Description:  "The origin server weight. Must be between `0` and `100`. Defaults to `0`. The higher the weight, the higher the priority. No effect when `type` is `code`.",
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"ttl": {
							Type:        schema.TypeInt,
							Description: "Cache duration in seconds.",
							Required:    true,
						},
					},
				},
			},
			"cache_url_args": {
				Type:        schema.TypeList,
				Description: "Cache parameter filter configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_full_url": {
							Type:        schema.TypeBool,
							Description: "Whether caching of full url is enabled. Defaults to `true`",
							Optional:    true,
							Default:     true,
						},
						"cache_url_args": {
							Type:        schema.TypeSet,
							Description: "Query parameters in url that will be preserved. No effect when `cache_full_url` is `true`.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: cacheUrlArgsInnerDiffSuppress,
						},
					},
				},
				DiffSuppressFunc: cacheUrlArgsDiffSuppress,
			},
			"error_page": {
				Type:        schema.TypeSet,
				Description: "Error page configuration of the acceleration domain.",
				Optional:    true,
				Set:         errorPageHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:         schema.TypeInt,
							Description:  "Error status code. Valid values: `401`, `403`, `404`, `405`, `414`, `429`, `500`, `501`, `502`, `503`, `504`.",
							Required:     true,
							ValidateFunc: validation.IntInSlice(validErrorPageStatusCodes()),
						},
						"url": {
							Type:        schema.TypeString,
							Description: "Destination address of redirection.",
							Required:    true,
						},
					},
				},
			},
			"cache_share": {
				Type:        schema.TypeList,
				Description: "Cache sharing configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether cache sharing is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"domain": {
							Type:        schema.TypeString,
							Description: "Another acceleration domain under current user. Must be set when enabled.",
							Optional:    true,
						},
					},
				},
				DiffSuppressFunc: cacheShareDiffSuppress,
			},
			"mobile_access": {
				Type:        schema.TypeList,
				Description: "Mobile access configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguish_client": {
							Type:        schema.TypeBool,
							Description: "Whether mobile access is enabled. Defaults to `false`",
							Optional:    true,
							Default:     false,
						},
					},
				},
				DiffSuppressFunc: mobileAccessDiffSuppress,
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
	return resourceDomainConfigCacheRead(d, meta)
}

func resourceDomainConfigCacheRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()
	d.Set("domain", domain)

	if err := readCommonConfigCache(d, conn, domain); err != nil {
		return err
	}
	if err := readCacheUrlArgs(d, conn, domain); err != nil {
		return err
	}
	if err := readMobileAccess(d, conn, domain); err != nil {
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

	return resourceDomainConfigCacheRead(d, meta)
}

func resourceDomainConfigCacheDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func readCommonConfigCache(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	config, err := FindDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config: %w", domain, err)
	}
	log.Printf("[DEBUG] Read CDN Domain (%s) Config CacheTTL result: %+v", domain, config.CacheTTL)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config ErrorPage result: %+v", domain, config.ErrorPage)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config CacheShare result: %+v", domain, config.CacheShare)

	d.Set("cache_ttl", flattenCacheTTLs(config.CacheTTL))
	d.Set("error_page", flattenErrorPages(config.ErrorPage))
	d.Set("cache_share", flattenCacheShare(config.CacheShare))

	return nil
}

func readCacheUrlArgs(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetCacheUrlArgs(domain)
	})
	log.Printf("[DEBUG] Read CDN Domain (%s) Config CacheUrlArgs result: %+v", domain, raw)

	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config CacheUrlArgs: %w", domain, err)
	}

	d.Set("cache_url_args", flattenCacheUrlArgs(raw.(*api.CacheUrlArgs)))
	return nil
}

func readMobileAccess(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetMobileAccess(domain)
	})
	log.Printf("[DEBUG] Read CDN Domain (%s) Config MobileAccess result: %+v", domain, raw)

	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config MobileAccess: %w", domain, err)
	}

	d.Set("mobile_access", flattenMobileAccess(raw.(bool)))
	return nil
}

func updateConfigCache(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateCacheTTL(d, conn, domain); err != nil {
		return err
	}
	if err := updateCacheUrlArgs(d, conn, domain); err != nil {
		return err
	}
	if err := updateErrorPage(d, conn, domain); err != nil {
		return err
	}
	if err := updateCacheShare(d, conn, domain); err != nil {
		return err
	}
	if err := updateMobileAccess(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func updateCacheTTL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("cache_ttl") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config CacheTTL", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			schema := d.Get("cache_ttl").(*schema.Set)
			return nil, client.SetCacheTTL(domain, expandCacheTTLs(schema))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config CacheTTL: %w", domain, err)
		}
	}
	return nil
}

func updateCacheUrlArgs(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("cache_url_args") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config CacheUrlArgs", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetCacheUrlArgs(domain, expandCacheUrlArgs(d.Get("cache_url_args").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config CacheUrlArgs: %w", domain, err)
		}
	}
	return nil
}

func updateErrorPage(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("error_page") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config ErrorPage", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetErrorPage(domain, expandErrorPages(d.Get("error_page").(*schema.Set)))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config ErrorPage: %w", domain, err)
		}
	}
	return nil
}

func updateCacheShare(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("cache_share") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config CacheShare", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			request := expandCacheShare(d.Get("cache_share").([]interface{}))
			return nil, client.SetCacheShared(domain, request)
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config CacheShare: %w", domain, err)
		}
	}
	return nil
}

func updateMobileAccess(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("mobile_access") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config MobileAccess", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetMobileAccess(domain, expandMobileAccess(d.Get("mobile_access").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config MobileAccess: %w", domain, err)
		}
	}
	return nil
}
