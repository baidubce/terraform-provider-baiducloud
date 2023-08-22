package iam

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func DataSourceAccessKeys() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query IAM access keys. \n\n",

		Read: dataSourceAccessKeysRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Description: "The name of the IAM user associated with the access keys.",
				Required:    true,
			},
			"access_keys": {
				Type:        schema.TypeList,
				Description: "The access key list.",
				Computed:    true,
				Elem:        AccessKeySchema(),
			},
		},
	}
}

func AccessKeySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Description: "The id of access key.",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the access key is enabled.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Date and time in RFC3339 format that the access key was created.",
				Computed:    true,
			},
			"last_used_time": {
				Type:        schema.TypeString,
				Description: "Date and time in RFC3339 format that the access key was last used.",
				Computed:    true,
			},
		},
	}
}

func dataSourceAccessKeysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	keyPairs, err := FindAccessKeys(conn, d.Get("username").(string))

	log.Printf("[DEBUG] Read IAM access keys result: %+v", keyPairs)
	if err != nil {
		return fmt.Errorf("error reading IAM access keys: %w", err)
	}

	if err := d.Set("access_keys", flattenAccessKeyList(keyPairs)); err != nil {
		return fmt.Errorf("error setting access_keys: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
