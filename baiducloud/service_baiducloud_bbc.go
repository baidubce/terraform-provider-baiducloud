package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

type BbcService struct {
	client *connectivity.BaiduClient
}

func (s *BbcService) ListAllBbcInstances(args *bbc.ListInstancesArgs) ([]bbc.InstanceModel, error) {
	action := "List all " + args.VpcId + " bbc instances"

	result := make([]bbc.InstanceModel, 0)
	for {
		raw, err := s.client.WithBbcClient(func(client *bbc.Client) (i interface{}, e error) {
			return client.ListInstances(args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*bbc.ListInstancesResult)
		result = append(result, response.Instances...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}
func (s *BbcService) GetBbcInstanceDetail(instanceId string) (*bbc.InstanceModel, error) {
	action := "Get " + instanceId + " bbc instances"

	for {
		raw, err := s.client.WithBbcClient(func(client *bbc.Client) (i interface{}, e error) {
			return client.GetInstanceDetail(instanceId)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*bbc.InstanceModel)
		return response, nil
	}
}
func (s *BbcService) FlattenBbcInstanceModelToMap(instances []bbc.InstanceModel) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(instances))

	for _, inst := range instances {
		result = append(result, map[string]interface{}{
			"instance_id":              inst.Id,
			"instance_name":            inst.Name,
			"hostname":                 inst.Hostname,
			"uuid":                     inst.Uuid,
			"description":              inst.Desc,
			"status":                   inst.Status,
			"payment_timing":           inst.PaymentTiming,
			"create_time":              inst.CreateTime,
			"expire_time":              inst.ExpireTime,
			"public_ip":                inst.PublicIp,
			"internal_ip":              inst.InternalIp,
			"rdma_ip":                  inst.RdmaIp,
			"image_id":                 inst.ImageId,
			"flavor_id":                inst.FlavorId,
			"zone":                     inst.Zone,
			"region":                   inst.Region,
			"has_alive":                inst.HasAlive,
			"tags":                     flattenTagsToMap(inst.Tags),
			"switch_id":                inst.SwitchId,
			"host_id":                  inst.HostId,
			"deployset_id":             inst.DeploysetId,
			"network_capacity_in_mbps": inst.NetworkCapacityInMbps,
			"rack_id":                  inst.RackId,
		})
	}

	return result, nil
}
func (s *BbcService) ListAllBbcImages(args *bbc.ListImageArgs) ([]bbc.ImageModel, error) {
	action := "List all " + args.ImageType + " bbc images"

	result := make([]bbc.ImageModel, 0)
	for {
		raw, err := s.client.WithBbcClient(func(client *bbc.Client) (i interface{}, e error) {
			return client.ListImage(args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*bbc.ListImageResult)
		result = append(result, response.Images...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}
func (s *BbcService) GetBbcImageDetails(imageId string) (*bbc.ImageModel, error) {
	action := "Get BBC image " + imageId + " detail"

	raw, err := s.client.WithBbcClient(func(client *bbc.Client) (i interface{}, e error) {
		return client.GetImageDetail(imageId)
	})
	if err != nil {
		return nil, WrapError(err)
	}
	addDebug(action, raw)

	response := raw.(*bbc.GetImageDetailResult)
	return response.Result, err
}
func (s *BbcService) FlattenImageModelToMap(images []bbc.ImageModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(images))

	for _, image := range images {
		result = append(result, map[string]interface{}{
			"id":          image.Id,
			"name":        image.Name,
			"type":        image.Type,
			"os_type":     image.OsType,
			"os_version":  image.OsVersion,
			"os_arch":     image.OsArch,
			"os_name":     image.OsName,
			"os_build":    image.OsBuild,
			"create_time": image.CreateTime,
			"status":      image.Status,
			"desc":        image.Desc,
		})
	}

	return result
}

func (s *BbcService) ListAllBbcFlavors() ([]bbc.FlavorModel, error) {
	action := "List all bbc flavors"

	result := make([]bbc.FlavorModel, 0)
	raw, err := s.client.WithBbcClient(func(client *bbc.Client) (i interface{}, e error) {
		return client.ListFlavors()
	})
	addDebug(action, raw)
	response := raw.(*bbc.ListFlavorsResult)
	result = append(result, response.Flavors...)
	return result, err
}
func (s *BbcService) FlattenFlavorModelToMap(flavors []bbc.FlavorModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(flavors))

	for _, flavor := range flavors {
		result = append(result, map[string]interface{}{
			"flavor_id":             flavor.FlavorId,
			"cpu_count":             flavor.CpuCount,
			"cpu_type":              flavor.CpuType,
			"memory_capacity_in_gb": flavor.MemoryCapacityInGB,
			"disk":                  flavor.Disk,
			"network_card":          flavor.NetworkCard,
			"others":                flavor.Others,
		})
	}
	return result
}
func (s *BbcService) getRaidIdByFlavor(flavorId string, raid string) (string, error) {
	action := "get Raid Id By Flavor"

	raw, err := s.client.WithBbcClient(func(client *bbc.Client) (i interface{}, e error) {
		return client.GetFlavorRaid(flavorId)
	})
	addDebug(action, raw)
	if err != nil {
		return "", err
	}
	response := raw.(*bbc.GetFlavorRaidResult)
	for _, elem := range response.Raids {
		if elem.Raid == raid {
			return elem.RaidId, nil
		}
	}
	return "", nil
}
func (s *BbcService) InstanceStateRefresh(instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query BBC instance " + instanceId
		raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return bbcClient.GetInstanceDetail(instanceId)
		})

		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(api.InstanceStatusDeleted), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}
		addDebug(action, raw)
		result, _ := raw.(*bbc.InstanceModel)
		return result, string(result.Status), nil
	}
}
func (s *BbcService) StartBbcInstance(instanceID string, timeout time.Duration) error {
	action := "Stop bbc instance " + instanceID

	_, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
		return nil, bbcClient.StartInstance(instanceID)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(bbc.InstanceStatusStopped), string(bbc.InstanceStatusStarting)},
		[]string{string(bbc.InstanceStatusRunning)},
		timeout,
		s.InstanceStateRefresh(instanceID),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	return nil
}
func (s *BbcService) StopBbcInstance(instanceID string, timeout time.Duration) error {
	action := "Stop bbc instance " + instanceID

	_, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
		return nil, bbcClient.StopInstance(instanceID, false)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(bbc.InstanceStatusStopping), string(bbc.InstanceStatusRunning)},
		[]string{string(bbc.InstanceStatusStopped)},
		timeout,
		s.InstanceStateRefresh(instanceID),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	return nil
}
func (*BbcService) updateBccInstanceDescription(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update Bcc Instance Description " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("description") {
		modifyInstanceDescArgs := &bbc.ModifyInstanceDescArgs{}
		modifyInstanceDescArgs.Description = d.Get("description").(string)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
				return nil, bbcClient.ModifyInstanceDesc(instanceID, modifyInstanceDescArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstanceDescArgs)
			return nil
		})
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}
		d.SetPartial("description")
	}
	return nil
}

func (*BbcService) updateBbcInstanceName(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update BBC Instance attribute " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("name") {
		modifyInstanceNameArgs := &bbc.ModifyInstanceNameArgs{}
		modifyInstanceNameArgs.Name = d.Get("name").(string)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
				return nil, bbcClient.ModifyInstanceName(instanceID, modifyInstanceNameArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstanceNameArgs)
			return nil
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("name")
	}

	return nil
}
func (*BbcService) updateBbcInstanceDesc(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update Instance description " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("description") {
		modifyInstanceDescArgs := &bbc.ModifyInstanceDescArgs{
			ClientToken: buildClientToken(),
		}
		modifyInstanceDescArgs.Description = d.Get("description").(string)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
				return nil, bbcClient.ModifyInstanceDesc(instanceID, modifyInstanceDescArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstanceDescArgs)
			return nil
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("description")
	}

	return nil
}
func (*BbcService) updateBbcInstanceAdminPass(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update Instance admin pass " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("admin_pass") {
		modifyInstancePasswordArgs := &bbc.ModifyInstancePasswordArgs{}
		modifyInstancePasswordArgs.AdminPass = d.Get("admin_pass").(string)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
				return nil, bbcClient.ModifyInstancePassword(instanceID, modifyInstancePasswordArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstancePasswordArgs)
			return nil
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("admin_pass")
	}

	return nil
}
func (s *BbcService) BbcImageStateRefresh(imageId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query custom image " + imageId
		raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return bbcClient.GetImageDetail(imageId)
		})

		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(bbc.ImageStatusNotAvailable), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_image", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}
		addDebug(action, raw)
		result, _ := raw.(*bbc.GetImageDetailResult)
		return result, string(result.Result.Status), nil
	}
}
func (*BbcService) updateBbcInstanceSecurityGroups(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update bbc instance security groups " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("security_groups") {
		o, n := d.GetChange("security_groups")

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		bindSGs := make([]interface{}, 0)
		groups, ok := d.GetOk("security_groups")
		if ok {
			bindSGs = groups.(*schema.Set).List()
		}
		unbindSGs := os.Difference(ns).List()
		// bind
		bindSecurityGroupIds := make([]string, 0)
		for _, sg := range bindSGs {
			bindSecurityGroupIds = append(bindSecurityGroupIds, sg.(string))
		}
		bindArgs := &bbc.BindSecurityGroupsArgs{
			// bind security groups
			InstanceIds:      append(make([]string, 0), instanceID),
			SecurityGroupIds: bindSecurityGroupIds,
		}
		if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return nil, bbcClient.BindSecurityGroups(bindArgs)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		// unbind
		for _, sg := range unbindSGs {
			// unbind security groups
			if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
				return nil, bbcClient.UnBindSecurityGroups(&bbc.UnBindSecurityGroupsArgs{
					InstanceId:      instanceID,
					SecurityGroupId: sg.(string),
				})
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
			}
		}

		d.SetPartial("security_groups")
	}

	return nil
}
func (s *BbcService) updateBbcInstanceAction(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update bbc instance action " + instanceID

	if d.HasChange("action") {
		act := d.Get("action").(string)
		addDebug(action, act)

		if act == INSTANCE_ACTION_START {
			if err := s.StartBbcInstance(instanceID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		} else if act == INSTANCE_ACTION_STOP {
			if err := s.StopBbcInstance(instanceID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}

		d.SetPartial("action")
	}

	return nil
}
