package cdn

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func ResourceDomainConfigAdvanced() *schema.Resource {

	return &schema.Resource{

		Description: "Use this resource to manage Advanced configuration of the acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN/s/Jjxzil1sd). \n\n" +
			"~> **NOTE:** Creating a resource will overwrite current Advanced configuration. " +
			"Deleting a resource won't change current configuration.",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceDomainConfigAdvancedCreate,
		Read:   resourceDomainConfigAdvancedRead,
		Update: resourceDomainConfigAdvancedUpdate,
		Delete: resourceDomainConfigAdvancedDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Name of the acceleration domain.",
				Required:    true,
				ForceNew:    true,
			},
			"ipv6_dispatch": {
				Type:        schema.TypeList,
				Description: "IPv6 configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Description: "Whether supporting for accessing cdn via ipv6 is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
					},
				},
				DiffSuppressFunc: ipv6DispatchDiffSuppress,
			},
			"http_header": {
				Type:        schema.TypeSet,
				Description: "Http header configuration of the acceleration domain.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Description:  "Header type. Valid values: `origin`(forwards to origin), `response`(responds to the user).",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"origin", "response"}, false),
						},
						"header": {
							Type:        schema.TypeString,
							Description: "Http header field. Can be a standard Header of HTTP, or a user-defined one, such as `x-bce-authorization`.",
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Value of the header field. A limited number of variables are supported.",
							Optional:    true,
						},
						"action": {
							Type:        schema.TypeString,
							Description: "Indicates whether to delete or add the header. Valid values: `add`, `remove`",
							Required:    true,
						},
						"describe": {
							Type:        schema.TypeString,
							Description: "Description. Length cannot exceed 100 characters.",
							Optional:    true,
						},
					},
				},
			},
			"media_drag": {
				Type:        schema.TypeList,
				Description: "Video dragging configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mp4": {
							Type:        schema.TypeList,
							Description: "Configuration of pseudo-streaming in mp4 type.",
							Optional:    true,
							MaxItems:    1,
							Elem:        schemaMediaCfg(true),
						},
						"flv": {
							Type:        schema.TypeList,
							Description: "Configuration of pseudo-streaming in flv type.",
							Optional:    true,
							MaxItems:    1,
							Elem:        schemaMediaCfg(false),
						},
					},
				},
			},
			"seo_switch": {
				Type:        schema.TypeList,
				Description: "SEO optimization configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directly_origin": {
							Type:         schema.TypeString,
							Description:  "Whether SEO optimization is enabled. Defaults to `OFF`. Other valid value: `ON`",
							Optional:     true,
							Default:      "OFF",
							ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
						},
					},
				},
				DiffSuppressFunc: seoSwitchDiffSuppress,
			},
			"file_trim": {
				Type:        schema.TypeBool,
				Description: "Whether page optimization is enabled. Defaults to `false`.",
				Optional:    true,
				Default:     false,
			},
			"compress": {
				Type:        schema.TypeList,
				Description: "Page compression configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:        schema.TypeBool,
							Description: "Whether page compression is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"type": {
							Type:         schema.TypeString,
							Description:  "Compression type. Defaults to `gzip`. Other valid values: `br`, `all`(br + gzip)",
							Optional:     true,
							Default:      "gzip",
							ValidateFunc: validation.StringInSlice([]string{"br", "gzip", "all"}, false),
						},
					},
				},
				DiffSuppressFunc: compressDiffSuppress,
			},
			"quic": {
				Type:        schema.TypeBool,
				Description: "Whether QUIC protocol is enabled. Defaults to `false`.",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func schemaMediaCfg(forMp4 bool) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_suffix": {
				Type:        schema.TypeSet,
				Description: "Video file extensions. For example, mp4 type may use [`mp4`, `m4v`, `m4a`], flv type may use `flv`.",
				Optional:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"start_arg_name": {
				Type:        schema.TypeString,
				Description: "Name of the start parameter. Defaults to `start`. Cannot be the same as `end_arg_name`.",
				Optional:    true,
				Default:     "start",
			},
			"end_arg_name": {
				Type:        schema.TypeString,
				Description: "Name of the start parameter. Defaults to `end`. Cannot be the same as `start_arg_name`.",
				Optional:    true,
				Default:     "end",
			},
			"drag_mode": {
				Type:         schema.TypeString,
				Description:  "Drag in seconds for mp4 type or in bytes for flv type. Valid value for mp4: `second`. Valid value for flv: `byteAV`, `byte`.",
				Required:     true,
				ValidateFunc: validation.StringInSlice(validDragMode(forMp4), false),
			},
		},
	}
}

func resourceDomainConfigAdvancedCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Get("domain").(string)

	if err := updateConfigAdvanced(d, conn, domain); err != nil {
		return err
	}

	d.SetId(domain)
	return resourceDomainConfigAdvancedRead(d, meta)
}

func resourceDomainConfigAdvancedRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()
	d.Set("domain", domain)

	if err := readCommonConfigAdvanced(d, conn, domain); err != nil {
		return err
	}
	if err := readIPv6Dispatch(d, conn, domain); err != nil {
		return err
	}
	if err := readCompress(d, conn, domain); err != nil {
		return err
	}

	return nil
}

func resourceDomainConfigAdvancedUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	if err := updateConfigAdvanced(d, conn, domain); err != nil {
		return err
	}
	return resourceDomainConfigAdvancedRead(d, meta)
}

func resourceDomainConfigAdvancedDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func readCommonConfigAdvanced(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	config, err := FindDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config: %w", domain, err)
	}
	log.Printf("[DEBUG] Read CDN Domain (%s) Config HttpHeader result: %+v", domain, config.HttpHeader)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config MediaDragConf result: %+v", domain, config.MediaDragConf)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config SeoSwitch result: %+v", domain, config.SeoSwitch)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config FileTrim result: %+v", domain, config.FileTrim)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config QUIC result: %+v", domain, config.QUIC)

	d.Set("http_header", flattenHttpHeaders(config.HttpHeader))
	d.Set("media_drag", flattenMediaDragConf(config.MediaDragConf))
	d.Set("seo_switch", flattenSeoSwitch(config.SeoSwitch))
	d.Set("file_trim", config.FileTrim)
	d.Set("quic", config.QUIC)
	return nil
}

func readIPv6Dispatch(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetIPv6(domain)
	})
	log.Printf("[DEBUG] Read CDN Domain (%s) Config IPv6Dispatch: %+v", domain, raw)

	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config IPv6Dispatch: %w", domain, err)
	}

	d.Set("ipv6_dispatch", flattenIPv6Dispatch(raw.(bool)))
	return nil
}

func readCompress(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetContentEncoding(domain)
	})
	log.Printf("[DEBUG] Read CDN Domain (%s) Config Compress result: %+v", domain, raw)

	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config Compress: %w", domain, err)
	}

	d.Set("compress", flattenCompress(raw.(string)))
	return nil
}

func updateConfigAdvanced(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateIPv6Dispatch(d, conn, domain); err != nil {
		return err
	}
	if err := updateHttpHeaderL(d, conn, domain); err != nil {
		return err
	}
	if err := updateMediaDrag(d, conn, domain); err != nil {
		return err
	}
	if err := updateSeoSwitch(d, conn, domain); err != nil {
		return err
	}
	if err := updateFileTrim(d, conn, domain); err != nil {
		return err
	}
	if err := updateCompress(d, conn, domain); err != nil {
		return err
	}
	if err := updateQUIC(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func updateIPv6Dispatch(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("ipv6_dispatch") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config IPv6Dispatch", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetIPv6(domain, expandIPv6Dispatch(d.Get("ipv6_dispatch").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config IPv6Dispatch: %w", domain, err)
		}
	}
	return nil
}

func updateHttpHeaderL(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("http_header") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config HttpHeader", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetHttpHeader(domain, expandHttpHeaders(d.Get("http_header").(*schema.Set)))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config HttpHeader: %w", domain, err)
		}
	}
	return nil
}

func updateMediaDrag(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("media_drag") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config MediaDrag", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetMediaDrag(domain, expandMediaDragConf(d.Get("media_drag").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config MediaDrag: %w", domain, err)
		}
	}
	return nil
}

func updateSeoSwitch(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("seo_switch") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config SeoSwitch", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetDomainSeo(domain, expandSeoSwitch(d.Get("seo_switch").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config SeoSwitch: %w", domain, err)
		}
	}
	return nil
}

func updateFileTrim(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("file_trim") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config FileTrim", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetFileTrim(domain, d.Get("file_trim").(bool))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config FileTrim: %w", domain, err)
		}
	}
	return nil
}

func updateCompress(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("compress") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config Compress", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			enabled, compressType := expandCompress(d.Get("compress").([]interface{}))
			return nil, client.SetContentEncoding(domain, enabled, compressType)
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config Compress: %w", domain, err)
		}
	}
	return nil
}

func updateQUIC(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("quic") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config QUIC", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetQUIC(domain, d.Get("quic").(bool))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config QUIC: %w", domain, err)
		}
	}
	return nil
}
