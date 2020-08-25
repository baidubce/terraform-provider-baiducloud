/*
Use this resource to create a CCEv2 Instance Group.

~> **NOTE:** The create/update/delete operation of ccev2 does NOT take effect immediately，maybe takes for several minutes.

Example Usage

```hcl
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_managed.id
    replicas = var.instance_group_replica_1
    instance_group_name = "ccev2_instance_group_1"
    instance_template {
      cce_instance_id = ""
      instance_name = "ccev2_test_instance"
      cluster_role = "node"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = var.vpc_id
        vpc_subnet_id = var.vpc_subnet_id
        security_group_id = var.security_group_id
        available_zone = "zoneA"
      }
      instance_resource {
        cpu = 4
        mem = 8
        root_disk_size = 40
        local_disk_size = 0
      }
      image_id = var.image_id
      instance_os {
        image_type = "System"
      }
      need_eip = false
      admin_password = "test123!YT"
      ssh_key_id = ""
      instance_charging_type = "Postpaid"
      runtime_type = "docker"
    }
  }
}
```
*/
package baiducloud

import (
	"errors"
	"log"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCCEv2InstanceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCCEv2InstanceGroupCreate,
		Read:   resourceBaiduCloudCCEv2InstanceGroupRead,
		Delete: resourceBaiduCloudCCEv2InstanceGroupDelete,
		Update: resourceBaiduCloudCCEv2InstanceGroupUpdate,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"spec": {
				Type:        schema.TypeList,
				Description: "Instance Group Spec",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Description: "Cluster ID of Instance Group",
							ForceNew:    true,
							Required:    true,
						},
						"instance_group_name": {
							Type:        schema.TypeString,
							Description: "Name of Instance Group",
							ForceNew:    true,
							Required:    true,
						},
						"instance_template": {
							Type:        schema.TypeList,
							Description: "Instance Spec of Instances in this Instance Group ",
							ForceNew:    true,
							Required:    true,
							MaxItems:    1,
							Elem:        resourceCCEv2InstanceSpec(),
						},
						"replicas": {
							Type:        schema.TypeInt,
							Description: "Number of instances in this Instance Group",
							Required:    true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeList,
				Description: "Instance Group Status",
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ready_replicas": {
							Type:        schema.TypeInt,
							Description: "Number of instances in RUNNING",
							Computed:    true,
						},
					},
				},
			},
			"nodes": {
				Type:        schema.TypeList,
				Description: "All detail info of nodes in this instance group",
				Computed:    true,
				Elem:        resourceCCEv2Instance(),
			},
		},
	}
}

func resourceBaiduCloudCCEv2InstanceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	arg, err := buildCreateInstanceGroupArgs(d)
	if err != nil {
		log.Printf("Build Create Instance Group Args Fail:" + err.Error())
		return WrapError(err)
	}

	action := "Create CCE Instance Group " + arg.Request.InstanceGroupName
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
			return client.CreateInstanceGroup(arg)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		resp := raw.(*ccev2.CreateInstanceGroupResponse)

		//waiting all instance in instance group are ready
		createTimeOutTime := d.Timeout(schema.TimeoutCreate)
		loopsCount := createTimeOutTime.Microseconds() / ((10 * time.Second).Microseconds())
		var i int64
		for i = 1; i <= loopsCount; i++ {
			time.Sleep(5 * time.Second)
			argsGetInstanceGroup := &ccev2.GetInstanceGroupArgs{
				ClusterID:       arg.ClusterID,
				InstanceGroupID: resp.InstanceGroupID,
			}
			rawInstanceGroupResp, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
				return client.GetInstanceGroup(argsGetInstanceGroup)
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			instanceGroupResp := rawInstanceGroupResp.(*ccev2.GetInstanceGroupResponse)
			if instanceGroupResp.InstanceGroup.Status.ReadyReplicas == instanceGroupResp.InstanceGroup.Spec.Replicas {
				break
			}
			if i == loopsCount {
				return resource.NonRetryableError(errors.New("create instance group time out"))
			}
		}
		addDebug(action, raw)
		response, ok := raw.(*ccev2.CreateInstanceGroupResponse)
		if !ok {
			err = errors.New("response format illegal")
			return resource.NonRetryableError(err)
		}
		d.SetId(response.InstanceGroupID)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCCEv2InstanceGroupRead(d, meta)
}

func resourceBaiduCloudCCEv2InstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	argsGetInstanceGroup, err := buildGetInstanceGroupArgs(d)
	if err != nil {
		log.Printf("Build Get Instance Group Args Fail:" + err.Error())
		return WrapError(err)
	}

	rawInstanceGroupResp, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.GetInstanceGroup(argsGetInstanceGroup)
	})

	if err != nil {
		return WrapError(err)
	}
	getInstanceGroupResp := rawInstanceGroupResp.(*ccev2.GetInstanceGroupResponse)

	if getInstanceGroupResp.InstanceGroup == nil || getInstanceGroupResp.InstanceGroup.Status == nil {
		return WrapError(errors.New("GetInstanceGroupResponse.InstanceGroup or  GetInstanceGroupResponse.InstanceGroup.Status is nil"))
	}

	statusList := make([]interface{}, 0)
	statusMap := make(map[string]interface{})
	statusMap["ready_replicas"] = getInstanceGroupResp.InstanceGroup.Status.ReadyReplicas
	statusList = append(statusList, statusMap)
	d.Set("status", statusList)

	argsGetInstanceOfInstanceGroup, err := buildGetInstancesOfInstanceGroupArgs(d)
	if err != nil {
		log.Printf("Build Get Instances of Instance Group Args Fail:" + err.Error())
		return WrapError(err)
	}
	rawInstancesResp, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.ListInstancesByInstanceGroupID(argsGetInstanceOfInstanceGroup)
	})
	if err != nil {
		return WrapError(err)
	}
	instancesResp := rawInstancesResp.(*ccev2.ListInstancesByInstanceGroupIDResponse)
	nodes, err := convertInstanceFromJsonToMap(instancesResp.Page.List, types.ClusterRoleNode)
	if err != nil {
		log.Printf("Get Instance Group Nodes Fail" + err.Error())
	} else {
		err = d.Set("nodes", nodes)
	}
	if err != nil {
		log.Printf("Set Instance Group Nodes Fail" + err.Error())
	}
	return nil
}

func resourceBaiduCloudCCEv2InstanceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args, err := buildUpdateInstanceGroupReplicaArgs(d)
	if err != nil {
		log.Printf("Build Update Instance Group Replica Args Fail:" + err.Error())
		return WrapError(err)
	}
	action := "Update CCE Instance Group: " + args.InstanceGroupID
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
			return client.UpdateInstanceGroupReplicas(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		//waiting all instance in instance group are ready
		createTimeOutTime := d.Timeout(schema.TimeoutCreate)
		loopsCount := createTimeOutTime.Microseconds() / ((10 * time.Second).Microseconds())
		var i int64
		for i = 1; i <= loopsCount; i++ {
			time.Sleep(5 * time.Second)
			argsGetInstanceGroup := &ccev2.GetInstanceGroupArgs{
				ClusterID:       args.ClusterID,
				InstanceGroupID: args.InstanceGroupID,
			}
			rawInstanceGroupResp, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
				return client.GetInstanceGroup(argsGetInstanceGroup)
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			instanceGroupResp := rawInstanceGroupResp.(*ccev2.GetInstanceGroupResponse)
			if instanceGroupResp.InstanceGroup.Status.ReadyReplicas == instanceGroupResp.InstanceGroup.Spec.Replicas {
				break
			}
			if i == loopsCount {
				return resource.NonRetryableError(errors.New("create instance group time out"))
			}
		}
		addDebug(action, raw)
		_, ok := raw.(*ccev2.UpdateInstanceGroupReplicasResponse)
		if !ok {
			err = errors.New("response format illegal")
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudCCEv2InstanceGroupRead(d, meta)
}

func resourceBaiduCloudCCEv2InstanceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	args, err := buildDeleteInstanceGroupArgs(d)
	if err != nil {
		log.Printf("Build Delete Instance Group Args Fail:" + err.Error())
		return WrapError(err)
	}
	action := "Delete CCE Instance Group: " + args.InstanceGroupID
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
			return client.DeleteInstanceGroup(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		time.Sleep(1 * time.Minute) //waiting for infrastructure delete before delete vpc & security group
		addDebug(action, raw)
		_, ok := raw.(*ccev2.DeleteInstanceGroupResponse)
		if !ok {
			err = errors.New("response format illegal")
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group", action, BCESDKGoERROR)
	}
	return nil
}
