package cdn

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func DataSourceDomainCertificate() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query the information of the certificate bound to the domain name. \n\n",

		Read: dataSourceDomainCertificateRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "Name of the acceleration domain.",
				Required:    true,
				ForceNew:    true,
			},
			"certificate": {
				Type:        schema.TypeList,
				Description: "Certificate information.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeString,
							Description: "Certificate ID.",
							Computed:    true,
						},
						"cert_name": {
							Type:        schema.TypeString,
							Description: "Certificate name.",
							Computed:    true,
						},
						"cert_common_name": {
							Type:        schema.TypeString,
							Description: "Common name of the certificate.",
							Computed:    true,
						},
						"cert_dns_names": {
							Type:        schema.TypeString,
							Description: "Other DNS name of Certificate.",
							Computed:    true,
						},
						"cert_start_time": {
							Type:        schema.TypeString,
							Description: "Effective time of certificate.",
							Computed:    true,
						},
						"cert_stop_time": {
							Type:        schema.TypeString,
							Description: "Expiration time of certificate.",
							Computed:    true,
						},
						"cert_create_time": {
							Type:        schema.TypeString,
							Description: "Creation time of certificate.",
							Computed:    true,
						},
						"cert_update_time": {
							Type:        schema.TypeString,
							Description: "Update time of certificate.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDomainCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	domain := d.Get("domain").(string)

	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetCert(domain)
	})
	log.Printf("[DEBUG] Read CDN Domain (%s) Certificate result: %+v", domain, raw)

	if err != nil {
		return fmt.Errorf("error getting CDN Domain (%s) certificate: %w", domain, err)
	}

	d.SetId(domain)
	d.Set("certificate", flattenDomainCertificate(raw.(*api.CertificateDetail)))

	return nil
}
