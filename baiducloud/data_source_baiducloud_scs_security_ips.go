/*
Use this data source to query SCS security ips.

Example Usage

```hcl
data "baiducloud_scs_security_ips" "default" {
	instance_id = "scs-xxxxx"
}

output "security_ips" {
 value = "${data.baiducloud_scs.default.security_ips}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudScsSecurityIps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudScsSecurityIpsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "ID of the instance",
				Required:    true,
				ForceNew:    true,
			},
			"security_ips": {
				Type:        schema.TypeList,
				Description: "security_ips",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Description: "securityIp",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudScsSecurityIpsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceID := d.Get("instance_id").(string)
	action := "Query SCS SecurityIp instanceID is " + instanceID

	raw, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
		return scsClient.GetSecurityIp(instanceID)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_security_ips", action, BCESDKGoERROR)
	}

	securityIpsResult, _ := raw.(*scs.GetSecurityIpResult)
	securityIps := make([]map[string]interface{}, 0)
	for _, ip := range securityIpsResult.SecurityIps {
		ipMap := make(map[string]interface{})
		ipMap["ip"] = ip
		securityIps = append(securityIps, ipMap)
	}
	addDebug(action, securityIps)

	FilterDataSourceResult(d, &securityIps)

	if err := d.Set("security_ips", securityIps); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_security_ips", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), securityIps); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_security_ips", action, BCESDKGoERROR)
		}
	}
	return nil
}
