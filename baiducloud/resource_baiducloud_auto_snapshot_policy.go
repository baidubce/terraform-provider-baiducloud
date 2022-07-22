/*
Provide a resource to create an AutoSnapshotPolicy.

Example Usage

```hcl
resource "baiducloud_auto_snapshot_policy" "my-asp" {
  name            = "${var.name}"
  time_points     = [0, 22]
  repeat_weekdays = [0, 3]
  retention_days  = -1
  volume_ids      = ["v-Trb3rQXa"]
}
```

Import

AutoSnapshotPolicy can be imported, e.g.

```hcl
$ terraform import baiducloud_auto_snapshot_policy.default id
```
*/
package baiducloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudAutoSnapshotPolicyCreate,
		Read:   resourceBaiduCloudAutoSnapshotPolicyRead,
		Update: resourceBaiduCloudAutoSnapshotPolicyUpdate,
		Delete: resourceBaiduCloudAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the automatic snapshot policy, which supports uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\",\".\", and the value must start with a letter, length 1-65.",
				Required:    true,
			},
			"time_points": {
				Type:        schema.TypeSet,
				Description: "Time point of generate snapshot in a day, the minimum unit is hour, supporting in range of [0, 23]",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(0, 23),
				},
			},
			"repeat_weekdays": {
				Type:        schema.TypeSet,
				Description: "Repeat time of the automatic snapshot policy, supporting in range of [0, 6]",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(0, 6),
				},
			},
			"retention_days": {
				Type:        schema.TypeInt,
				Description: "Number of days to retain the automatic snapshot, and -1 means permanently retained.",
				Required:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the automatic snapshot policy.",
				Computed:    true,
			},
			"created_time": {
				Type:        schema.TypeString,
				Description: "Creation time of the automatic snapshot policy.",
				Computed:    true,
			},
			"updated_time": {
				Type:        schema.TypeString,
				Description: "Update time of the automatic snapshot policy.",
				Computed:    true,
			},
			"deleted_time": {
				Type:        schema.TypeString,
				Description: "Deletion time of the automatic snapshot policy.",
				Computed:    true,
			},
			"last_execute_time": {
				Type:        schema.TypeString,
				Description: "Last execution time of the automatic snapshot policy.",
				Computed:    true,
			},
			"volume_ids": {
				Type:        schema.TypeSet,
				Description: "Volume id list to be attached of the automatic snapshot policy, these CDS volumes must be in-used.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"volume_count": {
				Type:        schema.TypeInt,
				Description: "The count of volumes associated with the snapshot.",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args, err := buildCreateAutoSnapshotPolicyArgs(d)
	if err != nil {
		return WrapError(err)
	}

	action := "Create Auto Snapshot Policy " + args.Name
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return client.CreateAutoSnapshotPolicy(args)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)

		response := raw.(*api.CreateASPResult)
		d.SetId(response.AspId)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policy", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudAutoSnapshotPolicyUpdate(d, meta)
}

func resourceBaiduCloudAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	id := d.Id()
	action := "Query Auto Snapshot Policy " + id

	raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
		return client.GetAutoSnapshotPolicy(id)
	})
	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policy", action, BCESDKGoERROR)
	}

	policy := &raw.(*api.GetASPDetailResult).AutoSnapshotPolicy
	d.Set("status", policy.Status)
	d.Set("created_time", policy.CreatedTime)
	d.Set("updated_time", policy.UpdatedTime)
	d.Set("deleted_time", policy.DeletedTime)
	d.Set("last_execute_time", policy.LastExecuteTime)
	d.Set("volume_count", policy.VolumeCount)
	d.Set("name", policy.Name)
	d.Set("repeat_weekdays", policy.RepeatWeekdays)
	d.Set("time_points", policy.TimePoints)
	d.Set("retention_days", policy.RetentionDays)

	return nil
}

func resourceBaiduCloudAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	id := d.Id()

	updatePolicy, updateArgs, err := buildUpdateAutoSnapshotPolicyArgs(d)
	if err != nil {
		return err
	}

	if updatePolicy {
		action := "Update Auto Snapshot Policy " + id

		updateArgs.AspId = id
		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return nil, client.UpdateAutoSnapshotPolicy(updateArgs)
		})

		if err != nil {
			if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policy", action, BCESDKGoERROR)
			}
		}
	}

	if d.HasChange("volume_ids") {
		o, n := d.GetChange("volume_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		if len(remove) > 0 {
			args := &api.DetachASPArgs{}
			for _, v := range remove {
				args.VolumeIds = append(args.VolumeIds, v.(string))
			}

			action := "Detach Auto Snapshot Policy " + id
			_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
				return nil, client.DetachAutoSnapshotPolicy(id, args)
			})

			if err != nil {
				if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
					return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policy", action, BCESDKGoERROR)
				}
			}
		}

		if len(add) > 0 {
			args := &api.AttachASPArgs{}
			for _, v := range add {
				args.VolumeIds = append(args.VolumeIds, v.(string))
			}

			action := "Attach Auto Snapshot Policy " + id
			_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
				return nil, client.AttachAutoSnapshotPolicy(id, args)
			})

			if err != nil {
				if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
					return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policy", action, BCESDKGoERROR)
				}
			}
		}
	}

	return resourceBaiduCloudAutoSnapshotPolicyRead(d, meta)
}

func resourceBaiduCloudAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	id := d.Id()
	action := "Delete Auto Snapshot Policy " + id

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return nil, client.DeleteAutoSnapshotPolicy(id)
		})
		addDebug(action, id)

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policy", action, BCESDKGoERROR)
	}

	return nil
}

func buildCreateAutoSnapshotPolicyArgs(d *schema.ResourceData) (*api.CreateASPArgs, error) {
	result := &api.CreateASPArgs{
		ClientToken: buildClientToken(),
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		result.Name = v.(string)
	}

	if v, ok := d.GetOk("time_points"); ok && v.(*schema.Set).Len() != 0 {
		for _, p := range v.(*schema.Set).List() {
			result.TimePoints = append(result.TimePoints, strconv.Itoa(p.(int)))
		}
	} else {
		return nil, fmt.Errorf("unset time_points")
	}

	if v, ok := d.GetOk("repeat_weekdays"); ok && v.(*schema.Set).Len() != 0 {
		for _, w := range v.(*schema.Set).List() {
			result.RepeatWeekdays = append(result.RepeatWeekdays, strconv.Itoa(w.(int)))
		}
	} else {
		return nil, fmt.Errorf("unset repeat_weekdyas")
	}

	if v, ok := d.GetOk("retention_days"); ok {
		result.RetentionDays = strconv.Itoa(v.(int))
	}

	return result, nil
}

func buildUpdateAutoSnapshotPolicyArgs(d *schema.ResourceData) (bool, *api.UpdateASPArgs, error) {
	result := &api.UpdateASPArgs{}
	update := false

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		result.Name = v.(string)

		if !update {
			update = d.HasChange("name")
		}
	}

	if v, ok := d.GetOk("time_points"); ok && v.(*schema.Set).Len() != 0 {
		for _, p := range v.(*schema.Set).List() {
			result.TimePoints = append(result.TimePoints, strconv.Itoa(p.(int)))
		}

		if !update {
			update = d.HasChange("time_points")
		}
	} else {
		return false, nil, fmt.Errorf("unset time_points")
	}

	if v, ok := d.GetOk("repeat_weekdays"); ok && v.(*schema.Set).Len() != 0 {
		for _, w := range v.(*schema.Set).List() {
			result.RepeatWeekdays = append(result.RepeatWeekdays, strconv.Itoa(w.(int)))
		}

		if !update {
			update = d.HasChange("repeat_weekdays")
		}
	} else {
		return false, nil, fmt.Errorf("unset repeat_weekdyas")
	}

	if v, ok := d.GetOk("retention_days"); ok {
		result.RetentionDays = strconv.Itoa(v.(int))

		if !update {
			update = d.HasChange("retention_days")
		}
	}

	return update, result, nil
}
