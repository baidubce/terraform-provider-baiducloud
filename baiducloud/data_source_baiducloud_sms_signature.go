/*
Use this data source to query sms signature .

Example Usage

```hcl
data "baiducloud_sms_signature" "default" {
	signature_id = "xxxxxx"
}

output "signature_info" {
 	value = "${data.baiducloud_sms_signature.default.signature_info}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudSMSSignature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudSMSSignatureRead,

		Schema: map[string]*schema.Schema{
			"signature_id": {
				Type:        schema.TypeString,
				Description: "signature id",
				Required:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Query result output file path",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"signature_info": {
				Type:        schema.TypeMap,
				Description: "signature content of sms",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Description: "signature content of sms",
							Computed:    true,
						},
						"content_type": {
							Type:        schema.TypeString,
							Description: "type of content",
							Computed:    true,
						},
						"country_type": {
							Type:        schema.TypeString,
							Description: "signature type of country",
							Computed:    true,
						},
						"user_id": {
							Type:        schema.TypeString,
							Description: "User id",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "status",
							Computed:    true,
						},
						"review": {
							Type:        schema.TypeString,
							Description: "commit review",
							Computed:    true,
						},
						"signature_id": {
							Type:        schema.TypeString,
							Description: "signature id",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudSMSSignatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	smsService := SMSService{client}
	var signatureId string
	if v, ok := d.GetOk("signature_id"); ok && v.(string) != "" {
		signatureId = v.(string)
	}

	action := "Query SMS signature " + signatureId
	signature, err := smsService.GetSMSSignatureDetail(signatureId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
	}
	signatureMap := map[string]interface{}{
		"content":      signature.Content,
		"content_type": signature.ContentType,
		"status":       signature.Status,
		"country_type": signature.CountryType,
		"review":       signature.Review,
		"user_id":      signature.UserId,
		"signature_id": signature.SignatureId,
	}

	if err := d.Set("signature_info", signatureMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), signatureMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
		}
	}

	return nil
}
