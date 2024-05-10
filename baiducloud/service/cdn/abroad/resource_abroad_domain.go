package abroad

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/model"
	"github.com/baidubce/bce-sdk-go/services/cdn/abroad"
	"github.com/baidubce/bce-sdk-go/services/cdn/abroad/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"

	"log"
)

func ResourceAbroadDomain() *schema.Resource {
	return &schema.Resource{

		Description: "Use this resource to manage abroad acceleration domain and its origin configuration. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/CDN-ABROAD/s/gjwvxiywx). \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceAbroadDomainCreate,
		Read:   resourceAbroadDomainRead,
		Update: resourceAbroadDomainUpdate,
		Delete: resourceAbroadDomainDelete,

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
						"type": {
							Type:         schema.TypeString,
							Description:  "origin type, value is IP or DOMAIN",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"IP", "DOMAIN"}, false),
						},
						"backup": {
							Type:        schema.TypeBool,
							Description: "Whether is a backup origin server. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"addr": {
							Type:        schema.TypeString,
							Description: "The ip address when forwarding to origin server",
							Required:    true,
						},
					},
				},
			},
			"designate_host_to_origin": {
				Type:         schema.TypeString,
				Description:  "Designate host to origin",
				Optional:     true,
			},
			"tags": flex.TagsSchema(),
		},
	}

}

func resourceAbroadDomainCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	domain := d.Get("domain").(string)

	_, err := conn.WithAbroadCdnClient(func(abroadCdnClient *abroad.Client) (interface{}, error) {
		tags := make([]model.TagModel, 0)
		var origin []api.OriginPeer
		if v, ok := d.GetOk("tags"); ok {
			tags = flex.TranceTagMapToModel(v.(map[string]interface{}))
		}
		if v, ok := d.GetOk("origin"); ok {
			origin = expandOriginPeers(v.([]interface{}))
		}
		log.Printf("[DEBUG] Create Abroad CDN Domain: %s %+v", domain, origin)
		return abroadCdnClient.CreateDomainWithOptions(domain, origin, abroad.CreateDomainWithTags(tags))
	})

	if err != nil {
		return fmt.Errorf("error creating Abroad CDN Domain (%s): %w", domain, err)
	}

	// set host to origin
	if err := updateDesignateHostToOrigin(d, conn, domain); err != nil {
		return err
	}

	d.SetId(domain)

	return resourceAbroadDomainRead(d, meta)
}

func resourceAbroadDomainRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	log.Printf("[DEBUG] Read Abroad CDN Domain (%s)", domain)
	config, err := FindAbroadDomainConfigByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error reading Abroad CDN Domain (%s): %w", domain, err)
	}
	err = d.Set("domain", config.Domain)
	if err != nil {
		return fmt.Errorf("error setting domain for Abroad CDN Domain (%s): %w", domain, err)
	}
	err = d.Set("origin", flattenAbroadOriginPeers(config.Origin))
	if err != nil {
		return fmt.Errorf("error setting origin for Abroad CDN Domain (%s): %w", domain, err)
	}
	tags, err := FindAbroadDomainTagsByName(conn, domain)
	if err != nil {
		return fmt.Errorf("error reading tags for Abroad CDN Domain (%s): %w", domain, err)
	}
	err = d.Set("tags", flex.FlattenTagsToMap(tags))
	if err != nil {
		return fmt.Errorf("error setting tag for Abroad CDN Domain (%s): %w", domain, err)
	}

	return nil
}

func resourceAbroadDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	if err := updateAbroadOrigins(d, conn, domain); err != nil {
		return err
	}
	if err := updateDesignateHostToOrigin(d, conn, domain); err != nil {
		return err
	}
	return nil
}

func resourceAbroadDomainDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Id()

	log.Printf("[DEBUG] Delete abroad CDN Domain: %s", domain)

	_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
		return nil, client.DeleteDomain(domain)
	})
	if err != nil {
		return fmt.Errorf("error deleting abroad CDN Domain (%s): %w", domain, err)
	}

	return nil
}

func expandOriginPeers(tfList []interface{}) []api.OriginPeer {
	var originPeers []api.OriginPeer
	for _, v := range tfList {
		tfMap := v.(map[string]interface{})
		originPeers = append(originPeers, api.OriginPeer{
			Type:   tfMap["type"].(string),
			Addr:   tfMap["addr"].(string),
			Backup: tfMap["backup"].(bool),
		})
	}
	return originPeers
}

func updateAbroadOrigins(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.HasChange("origin") {
		log.Printf("[DEBUG] Update Abroad CDN Domain origins(%s)", domain)

		origins := expandOriginPeers(d.Get("origin").([]interface{}))

		_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
			return nil, client.SetDomainOrigin(domain, origins)
		})
		if err != nil {
			return fmt.Errorf("error updating Abroad CDN Domain (%s) origins: %w", domain, err)
		}
	}
	return nil
}

func updateDesignateHostToOrigin(d *schema.ResourceData, conn *connectivity.BaiduClient, domain string) error {
	if d.IsNewResource() || d.HasChange("designate_host_to_origin") {
		log.Printf("[DEBUG] Update Abroad CDN Domain designate host to origin(%s)", domain)

		_, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
			return nil, client.SetHostToOrigin(domain, d.Get("designate_host_to_origin").(string))
		})
		if err != nil {
			return fmt.Errorf("error updating Abroad CDN Domain (%s) designate host to origin: %w", domain, err)
		}
	}
	return nil
}
