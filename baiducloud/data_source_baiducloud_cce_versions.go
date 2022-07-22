/*
Use this data source to list cce support kubernetes versions.

Example Usage

```hcl
data "baiducloud_cce_versions" "default" {}

output "versions" {
  value = "${data.baiducloud_cce_versions.default.versions}"
}
```
*/
package baiducloud

import (
	"regexp"

	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCceKubernetesVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCceKubernetesVersionRead,

		Schema: map[string]*schema.Schema{
			"version_regex": {
				Type:         schema.TypeString,
				Description:  "Regex pattern of the search version",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"versions": {
				Type:        schema.TypeList,
				Description: "Useful kubernetes version list",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceBaiduCloudCceKubernetesVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Query all cce versions"
	raw, err := client.WithCCEClient(func(client *cce.Client) (i interface{}, e error) {
		return client.ListVersions()
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_versions", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	var versionRegexStr string
	var versionRegex *regexp.Regexp
	if value, ok := d.GetOk("version_regex"); ok {
		versionRegexStr = value.(string)
		if len(versionRegexStr) > 0 {
			versionRegex = regexp.MustCompile(versionRegexStr)
		}
	}

	response := raw.(*cce.ListVersionsResult)
	versions := make([]string, 0, len(response.Data))
	for _, v := range response.Data {
		if len(versionRegexStr) > 0 && versionRegex != nil {
			if !versionRegex.MatchString(v) {
				continue
			}
		}
		versions = append(versions, v)
	}

	if err := d.Set("versions", versions); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_versions", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), versions); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_versions", action, BCESDKGoERROR)
		}
	}

	return nil
}
