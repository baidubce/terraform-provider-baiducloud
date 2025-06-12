/*
Use this resource to remove instances from a CCE InstanceGroup.

~> **NOTE:** After creation, it may take several minutes for the instances to be fully removed from the instance group.

Example Usage

```hcl
resource "baiducloud_ccev2_instance_group_detachment" "example" {
  cluster_id = "cce-example"
  instance_group_id = "cce-ig-example"
  instances_to_be_removed = ["cce-example-node"]
  clean_policy = "Delete"
  delete_option {
    move_out = false
    delete_resource = true
    delete_cds_snapshot = true
    drain_node = true
  }
}
```
*/
package baiducloud

import (
	"fmt"
	"time"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCCEv2InstanceGroupDetachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCCEv2InstanceGroupDetachmentCreate,
		Read:   resourceBaiduCloudCCEv2InstanceGroupDetachmentRead,
		Delete: resourceBaiduCloudCCEv2InstanceGroupDetachmentDelete,
		Update: resourceBaiduCloudCCEv2InstanceGroupDetachmentUpdate,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the CCE cluster.",
				ForceNew:    true,
				Required:    true,
			},
			"instance_group_id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance group.",
				ForceNew:    true,
				Required:    true,
			},
			"instances_to_be_removed": {
				Type:        schema.TypeSet,
				Description: "IDs of node to be removed. Note this refers to the node ID within the cluster, not the actual instance ID.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clean_policy": {
				Type:         schema.TypeString,
				Description:  "Whether to remove instances from the CCE cluster. `Remain` retains the instances in the cluster, `Delete` removes the instances from the cluster.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Remain", "Delete"}, false),
			},
			"delete_option": {
				Type:        schema.TypeList,
				Description: "Node deletion options.Required when `clean_policy` is set to `Delete`.",
				Optional:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"move_out": {
							Type:        schema.TypeBool,
							Description: "Whether to release the instance when removing the node. `true` keeps the instance, `false` releases it. Defaults to `true`.",
							Optional:    true,
							Default:     true,
						},
						"delete_resource": {
							Type:        schema.TypeBool,
							Description: "Whether to release related resources when removing the node. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"delete_cds_snapshot": {
							Type:        schema.TypeBool,
							Description: "Whether to delete associated CDS snapshots when removing the node. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"drain_node": {
							Type:        schema.TypeBool,
							Description: "Whether to perform node draining before removal. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, i interface{}) error {
			cleanPolicy := diff.Get("clean_policy").(string)
			if cleanPolicy == string(ccev2.CleanPolicyDelete) {
				if _, ok := diff.GetOk("delete_option"); !ok {
					return fmt.Errorf("'delete_option' must be set when 'clean_policy' is set to 'Delete'")
				}
			}
			return nil
		},
	}
}

func resourceBaiduCloudCCEv2InstanceGroupDetachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args := buildInstanceGroupDetachmentArgs(d)
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.CreateScaleDownInstanceGroupTask(args)
	})
	if err != nil {
		return err
	}

	resp := raw.(*ccev2.CreateTaskResp)
	d.SetId(resp.TaskID)

	if v, ok := d.GetOk("instances_to_be_removed"); ok && v.(*schema.Set).Len() > 0 && args.CleanPolicy == ccev2.CleanPolicyDelete {
		ccev2Service := Ccev2Service{client}
		instanceIds := expandStringSet(v.(*schema.Set))

		err := ccev2Service.waitForInstancesOperation([]string{EventStatusDeleting}, []string{EventStatusDeleted}, d.Timeout(schema.TimeoutCreate), instanceIds)
		if err != nil {
			return err
		}
	}

	return resourceBaiduCloudCCEv2InstanceGroupDetachmentRead(d, meta)
}

func resourceBaiduCloudCCEv2InstanceGroupDetachmentRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBaiduCloudCCEv2InstanceGroupDetachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBaiduCloudCCEv2InstanceGroupDetachmentDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildInstanceGroupDetachmentArgs(d *schema.ResourceData) *ccev2.CreateScaleDownInstanceGroupTaskArgs {
	args := &ccev2.CreateScaleDownInstanceGroupTaskArgs{}
	args.ClusterID = d.Get("cluster_id").(string)
	args.InstanceGroupID = d.Get("instance_group_id").(string)
	args.InstancesToBeRemoved = expandStringSet(d.Get("instances_to_be_removed").(*schema.Set))
	args.CleanPolicy = ccev2.CleanPolicy(d.Get("clean_policy").(string))

	if args.CleanPolicy == ccev2.CleanPolicyDelete {
		deleteOptionRaw := d.Get("delete_option").([]interface{})[0].(map[string]interface{})
		deleteOption := &types.DeleteOption{
			MoveOut:           deleteOptionRaw["move_out"].(bool),
			DeleteResource:    deleteOptionRaw["delete_resource"].(bool),
			DeleteCDSSnapshot: deleteOptionRaw["delete_cds_snapshot"].(bool),
			DrainNode:         deleteOptionRaw["drain_node"].(bool),
		}
		args.DeleteOption = deleteOption
	}

	return args
}
