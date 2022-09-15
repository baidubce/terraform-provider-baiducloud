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

func ResourceDomain() *schema.Resource {
	return &schema.Resource{

		Description: "Use this resource to manage acceleration domain and its origin configuration. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN/s/rjwvyev26). \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Name of the acceleration domain.",
				Required:    true,
			},
			"origin": {
				Type:        schema.TypeList,
				Description: "Origin server configuration of the acceleration domain.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer": {
							Type:        schema.TypeString,
							Description: "Format is {protocol://}{address}{:port}. `protocol` is optional, and valid values: `http`, `https`. `address` must be ip address or domain name. IPv6 address must be in '[ipv6]' format. `port` is optional.",
							Required:    true,
						},
						"host": {
							Type:        schema.TypeString,
							Description: "The host value used when forwarding to origin server",
							Optional:    true,
						},
						"backup": {
							Type:        schema.TypeBool,
							Description: "Whether is a backup origin server. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"weight": {
							Type:         schema.TypeInt,
							Description:  "The origin server weight. Must be between `1` and `100`. Sum of all weights should not be greater than 100. No effect when `peer` is domain name.",
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"isp": {
							Type:         schema.TypeString,
							Description:  "ISP of the origin server. Valid values: `un`(China Unicom), `ct`(China Telecom), `cm`(China Mobile)",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"un", "ct", "cm"}, false),
						},
					},
				},
			},
			"default_host": {
				Type:             schema.TypeString,
				Description:      "Domain-level host. Priority is lower than origin server host(ie origin.host).",
				Optional:         true,
				DiffSuppressFunc: defaultHostDiffSuppress,
			},
			"form": {
				Type:         schema.TypeString,
				Description:  "Business type of the domain name. Defaults to `default`. Valid values: `image`(small image file), `download`(large file downloading), `media` (streaming media on demand), `dynamic`(dynamic and static acceleration).",
				Optional:     true,
				Default:      "default",
				ValidateFunc: validation.StringInSlice(DomainFormValues(), false),
			},

			// computed
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the domain name. Possible values: `RUNNING`,`OPERATING`, `STOPPED`.",
				Computed:    true,
			},
			"cname": {
				Type:        schema.TypeString,
				Description: "CNAME address of the acceleration domain.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Creation time of the acceleration domain.",
				Computed:    true,
			},
			"last_modify_time": {
				Type:        schema.TypeString,
				Description: "Latest modification time of the acceleration domain.",
				Computed:    true,
			},
			"is_ban": {
				Type:        schema.TypeString,
				Description: "Whether the acceleration domain is blocked. `YES` means blocked, and `NO` means not blocked.",
				Computed:    true,
			},
		},
	}

}

func resourceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	domain := d.Get("domain").(string)
	input := buildCreateArgs(d)

	log.Printf("[DEBUG] Create CDN Domain: %s %+v", domain, input)

	_, err := conn.WithCdnClient(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.CreateDomain(domain, input)
	})

	if err != nil {
		return fmt.Errorf("error creating CDN Domain (%s): %w", domain, err)
	}

	d.SetId(domain)
	return resourceDomainRead(d, meta)
}

func resourceDomainRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	log.Printf("[DEBUG] Read CDN Domain (%s)", domain)

	config, err := FindDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error reading CDN Domain (%s): %w", domain, err)
	}

	d.Set("domain", domain)
	d.Set("form", config.Form)
	d.Set("origin", flattenOriginPeers(config.Origin))
	d.Set("default_host", config.DefaultHost)
	// computed
	d.Set("status", config.Status)
	d.Set("cname", config.Cname)
	d.Set("create_time", config.CreateTime)
	d.Set("last_modify_time", config.LastModifyTime)
	d.Set("is_ban", config.IsBan)

	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	if err := updateOrigins(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func resourceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	log.Printf("[DEBUG] Delete CDN Domain: %s", domain)

	_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return nil, client.DeleteDomain(domain)
	})
	if err != nil {
		return fmt.Errorf("error deleting CDN Domain (%s): %w", domain, err)
	}

	return nil
}

func buildCreateArgs(d *schema.ResourceData) *api.OriginInit {
	input := &api.OriginInit{}

	if v, ok := d.GetOk("origin"); ok {
		input.Origin = expandOriginPeers(v.([]interface{}))
	}
	if v, ok := d.GetOk("default_host"); ok {
		input.DefaultHost = v.(string)
	}

	if v, ok := d.GetOk("form"); ok {
		input.Form = v.(string)
	}

	return input
}

func updateOrigins(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.HasChanges("origin", "default_host") {
		log.Printf("[DEBUG] Update CDN Domain origins(%s)", domain)

		origins := expandOriginPeers(d.Get("origin").([]interface{}))
		defaultHost := d.Get("default_host").(string)
		if len(defaultHost) == 0 {
			defaultHost = domain
		}

		_, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
			return nil, client.SetDomainOrigin(domain, origins, defaultHost)
		})
		if err != nil {
			return fmt.Errorf("error updating CDN Domain (%s) origins: %w", domain, err)
		}
	}
	return nil
}
