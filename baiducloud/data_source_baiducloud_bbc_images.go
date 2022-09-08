/*
Use this data source to get bbc images list.

Example Usage

```hcl
data "baiducloud_bbc_images" "bbc_images" {
  image_type = "BbcSystem"
  os_name    = "CentOS"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBbcImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBbcImagesRead,

		Schema: map[string]*schema.Schema{
			"image_type": {
				Type:         schema.TypeString,
				Description:  "Image type of the images to be queried, support ALL/System/Custom/Integration/Sharing/GpuBccSystem/GpuBccCustom/FpgaBccSystem/FpgaBccCustom",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL", "System", "Custom", "Integration", "Sharing", "BbcSystem", "GpuBccCustom", "FpgaBccSystem", "FpgaBccCustom"}, false),
			},
			"name_regex": {
				Type:         schema.TypeString,
				Description:  "Regex pattern of the search image name",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"os_name": {
				Type:        schema.TypeString,
				Description: "Search image OS Name",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Images search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"images": {
				Type:        schema.TypeList,
				Description: "Image list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Image id",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Image name",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Image type",
							Computed:    true,
						},
						"os_type": {
							Type:        schema.TypeString,
							Description: "Image os type",
							Computed:    true,
						},
						"os_version": {
							Type:        schema.TypeString,
							Description: "Image os version",
							Computed:    true,
						},
						"os_arch": {
							Type:        schema.TypeString,
							Description: "Image os arch",
							Computed:    true,
						},
						"os_name": {
							Type:        schema.TypeString,
							Description: "Image os name",
							Computed:    true,
						},
						"os_build": {
							Type:        schema.TypeString,
							Description: "Image os build",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Image create time",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Image status",
							Computed:    true,
						},
						"desc": {
							Type:        schema.TypeString,
							Description: "Image description",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBbcImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	action := "Query All Bbc Images"
	listArgs := &bbc.ListImageArgs{}
	if v, ok := d.GetOk("image_type"); ok {
		listArgs.ImageType = v.(string)
	}

	imageResult, err := bbcService.ListAllBbcImages(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_images", action, BCESDKGoERROR)
	}
	addDebug(action, imageResult)

	imageListFilterName := make([]bbc.ImageModel, 0, len(imageResult))
	if v, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(v.(string))
		for _, image := range imageResult {
			if image.Name == "" {
				log.Printf("[WARN] Unable to find Image name to match against for image Id %q,"+
					"nothing to do.", image.Id)
				continue
			}

			if r.MatchString(image.Name) {
				imageListFilterName = append(imageListFilterName, image)
			}
		}
	} else {
		imageListFilterName = imageResult[:]
	}

	imageList := make([]bbc.ImageModel, 0, len(imageListFilterName))
	if v, ok := d.GetOk("os_name"); ok && v.(string) != "" {
		osName := strings.ToLower(v.(string))
		for _, image := range imageListFilterName {
			if image.OsName == "" {
				log.Printf("[WARN] Unable to find Image OS Name equal for image Id %q,"+
					"nothing to do.", image.Id)
				continue
			}

			if strings.ToLower(image.OsName) == osName {
				imageList = append(imageList, image)
			}
		}
	} else {
		imageList = imageListFilterName[:]
	}

	imageMap := bbcService.FlattenImageModelToMap(imageList)
	FilterDataSourceResult(d, &imageMap)

	if err := d.Set("images", imageMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_images", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), imageMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_images", action, BCESDKGoERROR)
		}
	}

	return nil
}
