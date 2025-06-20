package bec

import (
	"strconv"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/bec/api"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func flattenRegionList(regions []api.RegionInfo) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range regions {
		tfMap := map[string]interface{}{
			"region":       v.Region,
			"name":         v.Name,
			"country":      v.Country,
			"country_name": v.CountryName,
			"city_list":    flattenCityList(v.CityList),
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenCityList(cities []api.CityInfo) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range cities {
		tfMap := map[string]interface{}{
			"city":                  v.City,
			"name":                  v.Name,
			"service_provider_list": flattenServiceProviderList(v.ServiceProviderList),
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenServiceProviderList(serviceProviders []api.ServiceProviderInfo) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range serviceProviders {
		tfMap := map[string]interface{}{
			"service_provider": v.ServiceProvider,
			"name":             v.Name,
			"region_id":        v.RegionId,
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenVMInstances(vmInstances []api.VmInstanceDetailsVo) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range vmInstances {
		tfMap := map[string]interface{}{
			"service_id":          v.ServiceId,
			"vm_id":               v.VmId,
			"vm_name":             v.VmName,
			"host_name":           v.Hostname,
			"region_id":           v.RegionId,
			"spec":                v.Spec,
			"cpu":                 v.Cpu,
			"memory":              v.Mem,
			"image_type":          flattenImageType(v.OsImage.ImageType),
			"image_id":            v.OsImage.ImageId,
			"system_volume":       flattenSystemVolume(v.SystemVolume),
			"data_volume":         flattenDataVolumes(v.DataVolumeList),
			"need_public_ip":      v.NeedPublicIp,
			"need_ipv6_public_ip": v.NeedIpv6PublicIp,
			"bandwidth":           flattenBandwidth(v.Bandwidth),
			"dns_config":          flattenDNSConfig(v.Dns),
			"status":              v.Status,
			"internal_ip":         v.InternalIp,
			"public_ip":           v.PublicIp,
			"ipv6_public_ip":      v.Ipv6PublicIp,
			"create_time":         v.CreateTime,
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenImageType(imageType string) string {
	switch imageType {
	case "becCommon", "becCustom":
		return "bec"
	case "Custom":
		return "bcc"
	}
	return ""
}

func flattenBandwidth(bandwidth string) int {
	v, err := strconv.Atoi(strings.Trim(bandwidth, "Mbps"))
	if err != nil {
		return 0
	}
	return v
}

func flattenSystemVolume(systemVolume api.SystemVolumeConfig) interface{} {
	tfMap := map[string]interface{}{
		"name":        systemVolume.Name,
		"size_in_gb":  systemVolume.SizeInGB,
		"volume_type": string(systemVolume.VolumeType),
		"pvc_name":    systemVolume.PvcName,
	}
	return []interface{}{tfMap}
}

func expandSystemVolume(tfList []interface{}) *api.SystemVolumeConfig {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	tfMap := tfList[0].(map[string]interface{})
	return &api.SystemVolumeConfig{
		Name:       tfMap["name"].(string),
		SizeInGB:   tfMap["size_in_gb"].(int),
		VolumeType: api.DiskType(tfMap["volume_type"].(string)),
	}
}

func flattenDataVolumes(dataVolumes []api.VolumeConfig) interface{} {
	tfList := []interface{}{}
	for _, v := range dataVolumes {
		tfMap := map[string]interface{}{
			"name":        v.Name,
			"size_in_gb":  v.SizeInGB,
			"volume_type": string(v.VolumeType),
			"pvc_name":    v.PvcName,
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func expandDataVolumes(tfList []interface{}) *[]api.VolumeConfig {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	list := []api.VolumeConfig{}
	for _, v := range tfList {
		tfMap := v.(map[string]interface{})
		list = append(list, api.VolumeConfig{
			Name:       tfMap["name"].(string),
			SizeInGB:   tfMap["size_in_gb"].(int),
			VolumeType: api.DiskType(tfMap["volume_type"].(string)),
			PvcName:    tfMap["pvc_name"].(string),
		})
	}
	return &list
}

func flattenDNSConfig(dnsConfig string) interface{} {
	tfMap := map[string]interface{}{
		"dns_type": dnsTypeNone,
	}
	if len(dnsConfig) > 0 {
		if components := strings.Split(dnsConfig, "-"); len(components) > 0 {
			dnsType := strings.ToUpper(components[0])
			tfMap["dns_type"] = dnsType
			if dnsType == dnsTypeCustomize && len(components) > 1 {
				dnsAddress := strings.Split(components[1], ",")
				tfMap["dns_address"] = flex.FlattenStringValueList(dnsAddress)
			}
		}
	}
	return []interface{}{tfMap}
}

func expandDNSConfig(tfList []interface{}) *api.DnsConfig {
	dnsConfig := &api.DnsConfig{
		DnsType: dnsTypeNone,
	}
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	tfMap := tfList[0].(map[string]interface{})
	dnsConfig.DnsType = tfMap["dns_type"].(string)
	dnsConfig.DnsAddress = strings.Join(flex.ExpandStringValueList(tfMap["dns_address"].([]interface{})), ",")
	return dnsConfig
}

func flattenKeyConfig(keyPairs []api.KeyPair, adminPass string) interface{} {
	tfMap := map[string]interface{}{
		"type": "password",
	}
	if len(adminPass) > 0 {
		tfMap["admin_pass"] = adminPass
	}

	if len(keyPairs) > 0 {
		keyPairIds := []string{}
		for _, v := range keyPairs {
			keyPairIds = append(keyPairIds, v.KeyPairId)
		}
		tfMap["type"] = "bccKeyPair"
		tfMap["bcc_key_pair_id_list"] = flex.FlattenStringValueList(keyPairIds)
	}
	return []interface{}{tfMap}
}

func expandKeyConfig(tfList []interface{}) *api.KeyConfig {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	tfMap := tfList[0].(map[string]interface{})
	return &api.KeyConfig{
		Type:             tfMap["type"].(string),
		AdminPass:        tfMap["admin_pass"].(string),
		BccKeyPairIdList: flex.ExpandStringValueList(tfMap["bcc_key_pair_id_list"].([]interface{})),
	}
}

func expandNetworkConfigList(tfList []interface{}) *[]api.NetworkConfig {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	list := []api.NetworkConfig{}
	for _, v := range tfList {
		tfMap := v.(map[string]interface{})
		list = append(list, api.NetworkConfig{
			NodeType:     tfMap["node_type"].(string),
			NetworksList: expandNetworks(tfMap["networks"].([]interface{})),
		})
	}
	return &list
}

func expandNetworks(tfList []interface{}) *[]api.Networks {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	list := []api.Networks{}
	for _, v := range tfList {
		tfMap := v.(map[string]interface{})
		list = append(list, api.Networks{
			NetType: tfMap["net_type"].(string),
			NetName: tfMap["net_name"].(string),
		})
	}
	return &list
}
