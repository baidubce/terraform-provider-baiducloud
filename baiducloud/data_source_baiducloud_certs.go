/*
Use this data source to query CERT list.

Example Usage

```hcl
data "baiducloud_certs" "default" {
  name = "testCert"
}

output "certs" {
 value = "${data.baiducloud_certs.default.certs}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/cert"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCerts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCertsRead,

		Schema: map[string]*schema.Schema{
			"cert_name": {
				Type:        schema.TypeString,
				Description: "Name of the Cert to be queried",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Certs search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"certs": {
				Type:        schema.TypeList,
				Description: "A list of Cert",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeString,
							Description: "Cert's ID",
							Computed:    true,
						},
						"cert_type": {
							Type:        schema.TypeInt,
							Description: "Cert's type",
							Computed:    true,
						},
						"cert_name": {
							Type:        schema.TypeString,
							Description: "Cert's name",
							Computed:    true,
						},
						"cert_common_name": {
							Type:        schema.TypeString,
							Description: "Cert's common name",
							Computed:    true,
						},
						"cert_start_time": {
							Type:        schema.TypeString,
							Description: "Cert's start time",
							Computed:    true,
						},
						"cert_stop_time": {
							Type:        schema.TypeString,
							Description: "Cert's stop time",
							Computed:    true,
						},
						"cert_create_time": {
							Type:        schema.TypeString,
							Description: "Cert's create time",
							Computed:    true,
						},
						"cert_update_time": {
							Type:        schema.TypeString,
							Description: "Cert's update time",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudCertsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Query Certs"
	raw, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
		return client.ListCerts()
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_certs", action, BCESDKGoERROR)
	}
	addDebug(action, raw)
	response := raw.(*cert.ListCertResult)

	certName := ""
	if v, ok := d.GetOk("cert_name"); ok && v.(string) != "" {
		certName = v.(string)
	}
	certList := make([]map[string]interface{}, 0, len(response.Certs))
	for _, c := range response.Certs {
		if certName != "" && certName != c.CertName {
			continue
		}

		certList = append(certList, map[string]interface{}{
			"cert_id":          c.CertId,
			"cert_type":        c.CertType,
			"cert_name":        c.CertName,
			"cert_common_name": c.CertCommonName,
			"cert_start_time":  c.CertStartTime,
			"cert_stop_time":   c.CertStopTime,
			"cert_create_time": c.CertCreateTime,
			"cert_update_time": c.CertUpdateTime,
		})
	}

	FilterDataSourceResult(d, &certList)

	if err := d.Set("certs", certList); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_certs", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), certList); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_certs", action, BCESDKGoERROR)
		}
	}

	return nil
}
