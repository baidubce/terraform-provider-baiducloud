package cdn

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func ResourceDomainConfigHttps() *schema.Resource {

	return &schema.Resource{

		Description: "Use this resource to manage Https configuration of the acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN/s/Sjwvyf6w8). \n\n" +
			"~> **NOTE:** Creating a resource will overwrite current Https configuration. " +
			"Deleting a resource won't change current configuration.",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceDomainConfigHttpsCreate,
		Read:   resourceDomainConfigHttpsRead,
		Update: resourceDomainConfigHttpsUpdate,
		Delete: resourceDomainConfigHttpsDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Name of the acceleration domain.",
				Required:    true,
				ForceNew:    true,
			},
			"https": {
				Type:        schema.TypeList,
				Description: "Https configuration of the acceleration domain.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether https acceleration is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"cert_id": {
							Type:        schema.TypeString,
							Description: "SSL Certificate ID. Can be obtained through data source `baiducloud_cdn_domain_certificate`",
							Optional:    true,
						},
						"http_redirect": {
							Type:        schema.TypeBool,
							Description: "Whether redirecting HTTP requests to HTTPS is enabled. Defaults to `false`. Cannot be true at the same time as `https_redirect`",
							Optional:    true,
							Default:     false,
						},
						"http_redirect_code": {
							Type:         schema.TypeInt,
							Description:  "HTTP redirection status code. Defaults to `302`. Other valid value: `301`.",
							Optional:     true,
							Default:      302,
							ValidateFunc: validation.IntInSlice([]int{301, 302}),
						},
						"https_redirect": {
							Type:        schema.TypeBool,
							Description: "Whether redirecting HTTPS requests to HTTP is enabled. Defaults to `false`. Cannot be true at the same time as `http_redirect`",
							Optional:    true,
							Default:     false,
						},
						"https_redirect_code": {
							Type:         schema.TypeInt,
							Description:  "HTTPS redirection status code. Defaults to `302`. Other valid value: `301`.",
							Optional:     true,
							Default:      302,
							ValidateFunc: validation.IntInSlice([]int{301, 302}),
						},
						"http2_enabled": {
							Type:        schema.TypeBool,
							Description: "Whether HTTP2 feature is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"verify_client": {
							Type:        schema.TypeBool,
							Description: "Whether HTTPS two-way authentication is enabled. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"ssl_protocols": {
							Type:        schema.TypeSet,
							Description: "Supported TLS Versions. Valid values: `TLSv1.0`, `TLSv1.1`, `TLSv1.2`, `TLSv1.3`",
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"TLSv1.0", "TLSv1.1", "TLSv1.2", "TLSv1.3"}, false),
							},
						},
					},
				},
				DiffSuppressFunc: httpsDiffSuppress,
			},
			"ocsp": {
				Type:        schema.TypeBool,
				Description: "Whether OCSP stapling is enabled. Defaults to `false`.",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceDomainConfigHttpsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Get("domain").(string)

	if err := updateConfigHttps(d, conn, domain); err != nil {
		return err
	}

	d.SetId(domain)
	return resourceDomainConfigHttpsRead(d, meta)
}

func resourceDomainConfigHttpsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()
	d.Set("domain", domain)

	if err := readCommonConfigHttps(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func resourceDomainConfigHttpsUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	if err := updateConfigHttps(d, conn, domain); err != nil {
		return err
	}
	return resourceDomainConfigHttpsRead(d, meta)
}

func resourceDomainConfigHttpsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func readCommonConfigHttps(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	config, err := FindDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) Config: %w", domain, err)
	}
	log.Printf("[DEBUG] Read CDN Domain (%s) Config Https result: %+v", domain, config.Https)
	log.Printf("[DEBUG] Read CDN Domain (%s) Config OCSP result: %+v", domain, config.OCSP)

	d.Set("https", flattenHttps(config.Https))
	d.Set("ocsp", config.OCSP)
	return nil
}

func updateConfigHttps(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateHttps(d, conn, domain); err != nil {
		return err
	}
	if err := updateOCSP(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func updateHttps(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("https") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config Https", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetDomainHttps(domain, expandHttps(d.Get("https").([]interface{})))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config Https: %w", domain, err)
		}
	}
	return nil
}

func updateOCSP(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("ocsp") {
		log.Printf("[DEBUG] Update CDN Domain (%s) Config OCSP", domain)

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetOCSP(domain, d.Get("ocsp").(bool))
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) Config OCSP: %w", domain, err)
		}
	}
	return nil
}
