package eip

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceEipGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to attach one or more EIPs to an EIP Group. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/EIP/s/ukoslf7lm). \n\n",

		Create: resourceEipGroupAttachmentCreate,
		Read:   flex.DoNothing,
		Delete: flex.DoNothing,

		Schema: map[string]*schema.Schema{
			"eip_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the EIP Group.",
			},
			"eips": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of EIPs to be attached, including both IPv4 and IPv6 addresses.",
			},
		},
	}
}

func resourceEipGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	eipGroupID := d.Get("eip_group_id").(string)
	eips := flex.ExpandStringValueSet(d.Get("eips").(*schema.Set))

	_, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
		args := eip.EipGroupMoveInArgs{Eips: eips}
		return nil, client.EipGroupMoveIn(eipGroupID, &args)
	})
	if err != nil {
		return fmt.Errorf("error attaching eips to eip group (%s): %w", eipGroupID, err)
	}

	d.SetId(resource.UniqueId())
	return nil

}
