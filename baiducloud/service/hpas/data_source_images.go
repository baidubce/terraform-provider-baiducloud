package hpas

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func DataSourceImages() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query HPAS image list. \n\n",

		Read: dataSourceImagesRead,

		Schema: map[string]*schema.Schema{
			"image_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
				Description:  "Image type. Valid values: `System`, `Custom`.",
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application type that the image supports. e.g., `llama2_7B_train`.",
			},
			"image_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        ImageSchema(),
				Description: "Image list.",
			},
		},
	}
}

func ImageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The short ID of the image.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the image.",
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the image. Possible values: `System`, `Custom`.",
			},
			"image_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the image. Possible values: `Creating`, `CreatedFailed`, `Available`, `NotAvailable`, `Error`.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the image.",
			},
			"supported_app_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of supported application types.",
			},
		},
	}
}

func dataSourceImagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	result, err := FindImages(conn, d.Get("image_type").(string))
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("app_type"); ok {
		filtered := []api.ImageResponse{}
		for _, image := range result {
			for _, appType := range image.SupportedAppType {
				if appType == v {
					filtered = append(filtered, image)
					break
				}
			}
		}
		result = filtered
	}

	if err := d.Set("image_list", flattenImageList(result)); err != nil {
		return fmt.Errorf("error setting image_list: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
