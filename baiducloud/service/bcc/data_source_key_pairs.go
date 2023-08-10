package bcc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func DataSourceKeyPairs() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query BCC key pairs. \n\n",

		Read: dataSourceKeyPairsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of key pair. Use this to filter key pair list.",
				Optional:    true,
			},
			"key_pairs": {
				Type:        schema.TypeList,
				Description: "The key pair list.",
				Computed:    true,
				Elem:        KeyPairSchema(),
			},
		},
	}
}

func KeyPairSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"keypair_id": {
				Type:        schema.TypeString,
				Description: "The id of key pair.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of key pair.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of key pair.",
				Computed:    true,
			},
			"created_time": {
				Type:        schema.TypeString,
				Description: "The creation time of key pair.",
				Computed:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "The public key of keypair.",
				Computed:    true,
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Description: "The number of instances bound to key pair.",
				Computed:    true,
			},
			"region_id": {
				Type:        schema.TypeString,
				Description: "The id of the region to which key pair belongs.",
				Computed:    true,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Description: "The fingerprint of key pair.",
				Computed:    true,
			},
		},
	}
}

func dataSourceKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	keyPairs, err := FindKeyPairs(conn, d.Get("name").(string))

	log.Printf("[DEBUG] Read BCC key pairs result: %+v", keyPairs)
	if err != nil {
		return fmt.Errorf("error reading BCC key pairs: %w", err)
	}

	if err := d.Set("key_pairs", flattenKeyPairList(keyPairs)); err != nil {
		return fmt.Errorf("error setting key_pairs: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
