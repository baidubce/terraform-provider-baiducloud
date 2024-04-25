package mongodb

import "github.com/baidubce/bce-sdk-go/services/mongodb"

func flattenTags(tags []mongodb.TagModel) map[string]string {
	if tags == nil || len(tags) == 0 {
		return nil
	}
	tfMap := map[string]string{}
	for _, v := range tags {
		tfMap[v.TagKey] = v.TagValue
	}
	return tfMap
}

func expandTags(tfMap map[string]interface{}) []mongodb.TagModel {
	if tfMap == nil || len(tfMap) == 0 {
		return nil
	}
	var tagList []mongodb.TagModel
	for k, v := range tfMap {
		tagList = append(tagList, mongodb.TagModel{
			TagKey:   k,
			TagValue: v.(string),
		})
	}
	return tagList
}

func flattenSubnets(subnets []mongodb.SubnetMap) interface{} {
	if subnets == nil || len(subnets) == 0 {
		return nil
	}
	tfList := make([]map[string]interface{}, 0)
	for _, v := range subnets {
		tfMap := map[string]interface{}{
			"subnet_id": v.SubnetId,
			"zone_name": v.ZoneName,
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func expandSubnets(tfList []interface{}) []mongodb.SubnetMap {
	if tfList == nil || len(tfList) == 0 {
		return nil
	}
	var subnetList []mongodb.SubnetMap
	for _, v := range tfList {
		item := v.(map[string]interface{})
		subnetList = append(subnetList, mongodb.SubnetMap{
			SubnetId: item["subnet_id"].(string),
			ZoneName: item["zone_name"].(string),
		})
	}
	return subnetList
}

func flattenInstanceList(instanceList []mongodb.InstanceModel) interface{} {
	tfList := make([]map[string]interface{}, 0)
	for _, v := range instanceList {
		tfMap := map[string]interface{}{
			"instance_id":       v.DbInstanceId,
			"name":              v.DbInstanceName,
			"payment_timing":    v.PaymentTiming,
			"vpc_id":            v.VpcId,
			"subnets":           flattenSubnets(v.Subnets),
			"tags":              flattenTags(v.Tags),
			"type":              v.DbInstanceType,
			"storage_engine":    v.StorageEngine,
			"engine_version":    v.EngineVersion,
			"status":            v.DbInstanceStatus,
			"connection_string": v.ConnectionString,
			"create_time":       v.CreateTime.Local().Format("2006-01-02 15:04:05"),
			"cpu_count":         v.DbInstanceCpuCount,
			"memory_capacity":   v.DbInstanceMemoryCapacity,
			"storage":           v.DbInstanceStorage,
			"voting_member_num": v.VotingMemberNum,
			"readonly_node_num": v.ReadonlyNodeNum,
			"port":              v.Port,
			"mongos_count":      v.MongosCount,
			"shard_count":       v.ShardCount,
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenNodeList(nodeList []mongodb.NodeModel) interface{} {
	tfList := make([]map[string]interface{}, 0)
	for _, v := range nodeList {
		tfMap := map[string]interface{}{
			"node_id":           v.NodeId,
			"name":              v.Name,
			"status":            v.Status,
			"cpu_count":         v.CpuCount,
			"memory_capacity":   v.MemoryCapacity,
			"storage":           v.Storage,
			"storage_type":      v.StorageType,
			"connection_string": v.ConnectString,
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}
