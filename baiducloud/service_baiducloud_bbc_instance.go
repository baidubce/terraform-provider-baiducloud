package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/hashicorp/terraform/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type BbcService struct {
	client *connectivity.BaiduClient
}

func (s *BbcService) InstanceBbcStateRefresh(instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query BBC instance " + instanceId
		raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return bbcClient.GetInstanceDetail(instanceId)
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(bbc.InstanceStatusDeleted), nil
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

func (s *BbcService) ListAllInstance(args *bbc.ListInstancesArgs) ([]bbc.InstanceModel, error) {
	result := make([]bbc.InstanceModel, 0)

	action := "List all BBC instance "
	for {
		raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return bbcClient.ListInstances(args)
		})
		if err != nil {
			return nil, err
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

func (s *BbcService) GetInstanceDetail(instanceID string) (*bbc.InstanceModel, error) {
	action := "Get instance detail " + instanceID

	raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
		return bbcClient.GetInstanceDetail(instanceID)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	result, _ := raw.(*bbc.InstanceModel)
	return result, nil
}

func (s *BbcService) GetSystemVolume(instanceId string) (*bbc.VolumeModel, error) {
	action := "Get system volume " + instanceId

	volList, err := s.ListAllVolumes(instanceId)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	for _, vol := range volList {
		if vol.Type == bbc.VolumeTypeSYSTEM {
			return &vol, nil
		}
	}

	return nil, nil
}
func (s *BbcService) GetFlavorDetail(flavorId string) (*bbc.GetFlavorDetailResult, error) {
	raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
		return bbcClient.GetFlavorDetail(flavorId)
	})
	return raw.(*bbc.GetFlavorDetailResult), err
}

func (s *BbcService) GetFlavors() (*bbc.ListFlavorsResult, error) {
	raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
		return bbcClient.ListFlavors()
	})
	return raw.(*bbc.ListFlavorsResult), err
}

func (s *BbcService) GetRaids(flavorId string) (*bbc.GetFlavorRaidResult, error) {
	raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
		return bbcClient.GetFlavorRaid(flavorId)
	})
	return raw.(*bbc.GetFlavorRaidResult), err
}
func (s *BbcService) FlattenFlavorsToMap(flavors []bbc.FlavorModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(flavors))

	for _, item := range flavors {
		result = append(result, map[string]interface{}{
			"flavor_id":             item.FlavorId,
			"cpu_count":             item.CpuCount,
			"disk":                  item.Disk,
			"memory_capacity_in_gb": item.MemoryCapacityInGB,
			"network_card":          item.NetworkCard,
		})
	}
	return result
}

func (s *BbcService) FlattenRaidsToMap(raids []bbc.RaidModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(raids))

	for _, item := range raids {
		result = append(result, map[string]interface{}{
			"raid_id":        item.RaidId,
			"raid":           item.Raid,
			"sys_swap_size":  item.SysSwapSize,
			"sys_root_size":  item.SysRootSize,
			"sys_home_size":  item.SysHomeSize,
			"sys_disk_size":  item.SysDiskSize,
			"data_disk_size": item.DataDiskSize,
		})
	}
	return result
}
func (s *BbcService) FlattenFlavorDetailToMap(flavor *bbc.GetFlavorDetailResult) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	result = append(result, map[string]interface{}{
		"flavor_id":             flavor.FlavorId,
		"cpu_count":             flavor.CpuCount,
		"disk":                  flavor.Disk,
		"memory_capacity_in_gb": flavor.MemoryCapacityInGB,
		"network_card":          flavor.NetworkCard,
	})
	return result
}

func (s *BbcService) ListAllVolumes(instanceId string) ([]bbc.VolumeModel, error) {
	args := &bbc.ListCDSVolumeArgs{
		MaxKeys:    100,
		InstanceId: instanceId,
	}
	action := "List all volumes " + instanceId

	volList := make([]bbc.VolumeModel, 0)
	for {
		raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return bbcClient.ListCDSVolume(args)
		})
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}
		addDebug(action, raw)

		result, _ := raw.(*bbc.ListCDSVolumeResult)
		volList = append(volList, result.Volumes...)

		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
	}

	return volList, nil
}

func (s *BbcService) FlattenInstanceModelToMap(instances []bbc.InstanceModel) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(instances))

	for _, inst := range instances {
		flavor, _ := s.GetFlavorDetail(inst.FlavorId)
		args := &bbc.GetVpcSubnetArgs{
			BbcIds: []string{inst.Id},
		}
		raw, err := s.client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
			return bbcClient.GetVpcSubnet(args)
		})

		subnet_id := ""
		vpc_id := ""

		if err == nil {
			vpcs, _ := raw.(*bbc.GetVpcSubnetResult)
			for _, sg := range vpcs.NetworkInfo {
				subnet_id = sg.Subnet.SubnetId
				vpc_id = sg.Vpc.VpcId
			}
		}

		volumes, err := s.ListAllVolumes(inst.Id)
		cdsDisks := make([]interface{}, 0, len(volumes))
		if err == nil {
			for _, vol := range volumes {
				cdsMap := make(map[string]interface{})
				cdsMap["disk_size_in_gb"] = vol.DiskSizeInGB
				cdsMap["storage_type"] = vol.StorageType
				cdsMap["is_system_volume"] = vol.IsSystemVolume
				cdsDisks = append(cdsDisks, cdsMap)
			}
		}

		result = append(result, map[string]interface{}{
			"instance_id":              inst.Id,
			"flavor_id":                inst.FlavorId,
			"name":                     inst.Name,
			"cpu_count":                flavor.CpuCount,
			"memory_capacity_in_gb":    flavor.MemoryCapacityInGB,
			"status":                   inst.Status,
			"description":              inst.Desc,
			"payment_timing":           inst.PaymentTiming,
			"create_time":              inst.CreateTime,
			"expire_time":              inst.ExpireTime,
			"internal_ip":              inst.InternalIp,
			"public_ip":                inst.PublicIp,
			"image_id":                 inst.ImageId,
			"network_capacity_in_mbps": inst.NetworkCapacityInMbps,
			"zone_name":                inst.Zone,
			"subnet_id":                subnet_id,
			"vpc_id":                   vpc_id,
			"cds_disks":                cdsDisks,
			"tags":                     flattenTagsToMap(inst.Tags),
		})
	}

	return result, nil
}
