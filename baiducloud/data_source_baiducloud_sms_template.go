/*
Use this data source to query sms template .

Example Usage

```hcl
data "baiducloud_sms_template" "default" {
	template_id = "xxxxxx"
}

output "template_info" {
 	value = "${data.baiducloud_sms_template.default.template_info}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudSMSTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudSMSTemplateRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Description: "template id",
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
			"template_info": {
				Type:        schema.TypeMap,
				Description: "template content of sms",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Template name of sms",
							Required:    true,
							ForceNew:    true,
						},
						"content": {
							Type:        schema.TypeString,
							Description: "Template content of sms",
							Required:    true,
							ForceNew:    true,
						},
						"sms_type": {
							Type:        schema.TypeString,
							Description: "Type of sms",
							Required:    true,
							ForceNew:    true,
						},
						"country_type": {
							Type:        schema.TypeString,
							Description: "Template type of country",
							Required:    true,
							ForceNew:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description of template",
							Required:    true,
							ForceNew:    true,
						},
						"template_id": {
							Type:        schema.TypeString,
							Description: "Template id",
							Optional:    true,
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
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudSMSTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	smsService := SMSService{client}
	var templateId string
	if v, ok := d.GetOk("template_id"); ok && v.(string) != "" {
		templateId = v.(string)
	}

	action := "Query SMS template " + templateId
	template, err := smsService.GetSMSTemplateDetail(templateId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_template", action, BCESDKGoERROR)
	}
	templateMap := map[string]interface{}{
		"name":         template.Name,
		"content":      template.Content,
		"country_type": template.CountryType,
		"sms_type":     template.SmsType,
		"description":  template.Description,
		"status":       template.Status,
		"review":       template.Review,
		"user_id":      template.UserId,
		"template_id":  template.TemplateId,
	}

	if err := d.Set("template_info", templateMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_template", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), templateMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_template", action, BCESDKGoERROR)
		}
	}

	return nil
}
