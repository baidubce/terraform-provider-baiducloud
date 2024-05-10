package abroad

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/cdn/abroad"
	"github.com/baidubce/bce-sdk-go/services/cdn/abroad/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"

	"log"
)

func ResourceAbroadDomainConfigHttps() *schema.Resource {
	return &schema.Resource{

		Description: "Use this resource to manage Https configuration of the abroad acceleration domain. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN-ABROAD/s/ckb0fx9ea). \n\n" +
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
				Type: schema.TypeBool,
				Description: "Redirect HTTP requests to HTTPS (redirect status code is 301). Defaults to `false`." +
					" This item is invalid when enabled=false. ",
				Optional: true,
				Default:  false,
			},
			"http2_enabled": {
				Type: schema.TypeBool,
				Description: "Whether HTTP2 feature is enabled. Defaults to `true`. " +
					"This item is invalid when enabled=false. ",
				Optional: true,
				Default:  true,
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

func updateConfigHttps(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if err := updateHttps(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func updateHttps(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	log.Printf("[DEBUG] Update abroad CDN Domain (%s) Config Https", domain)
	var options []api.HTTPSConfigOption
	options = append(options, api.HTTPSConfigCertID(d.Get("cert_id").(string)))
	if v, ok := d.GetOk("http_redirect"); ok {
		if v.(bool) {
			options = append(options, api.HTTPSConfigRedirectWith301())
		}
	}
	if v, ok := d.GetOk("http2_enabled"); ok {
		if v.(bool) {
			options = append(options, api.HTTPSConfigEnableH2())
		}
	}
	enabled := d.Get("enabled").(bool)
	_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
		return nil, client.SetHTTPSConfigWithOptions(domain, enabled, options...)
	})
	if err != nil {
		return fmt.Errorf("error updating abroad CDN Domain (%s) Config Https: %w", domain, err)
	}
	return nil
}

func resourceDomainConfigHttpsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()
	if err := d.Set("domain", domain); err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) domain: %w", domain, err)
	}

	domainConfig, _ := FindAbroadDomainConfigByName(conn, domain)
	if err := d.Set("enabled", domainConfig.HTTPSConfig.Enabled); err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) enabled: %w", domain, err)
	}
	if err := d.Set("cert_id", domainConfig.HTTPSConfig.CertId); err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) cert_id: %w", domain, err)
	}
	if err := d.Set("http_redirect", domainConfig.HTTPSConfig.HttpRedirect); err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) http_redirect: %w", domain, err)
	}
	if err := d.Set("http2_enabled", domainConfig.HTTPSConfig.Http2Enabled); err != nil {
		return fmt.Errorf("error reading abroad CDN Domain (%s) http2_enabled: %w", domain, err)
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
