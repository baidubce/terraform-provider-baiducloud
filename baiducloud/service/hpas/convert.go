package hpas

import (
	bcc "github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/baidubce/bce-sdk-go/services/hpas"
	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func flattenImageList(images []api.ImageResponse) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range images {
		tfMap := map[string]interface{}{
			"image_id":           v.ImageId,
			"name":               v.Name,
			"image_type":         v.ImageType,
			"image_status":       v.ImageStatus,
			"create_time":        v.CreateTime,
			"supported_app_type": flex.FlattenStringValueList(v.SupportedAppType),
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenInstanceList(instances []api.HpasResponse) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range instances {
		tfMap := map[string]interface{}{
			"payment_timing":        v.ChargeType,
			"tags":                  flex.FlattenTagModelToMap(v.Tags),
			"hpas_id":               v.HpasId,
			"app_type":              v.AppType,
			"app_performance_level": v.AppPerformanceLevel,
			"name":                  v.Name,
			"zone_name":             v.ZoneName,
			"image_id":              v.ImageId,
			"image_name":            v.ImageName,
			"internal_ip":           v.InternalIp,
			"subnet_id":             v.SubnetId,
			"subnet_name":           v.SubnetName,
			"vpc_id":                v.VpcId,
			"vpc_name":              v.VpcName,
			"vpc_cidr":              v.VpcCidr,
			"ehc_cluster_id":        v.EhcClusterId,
			"ehc_cluster_name":      v.EhcClusterName,
			"status":                v.Status,
			"create_time":           v.CreateTime,
		}
		if len(v.NicInfo) > 0 {
			tfMap["security_group_type"] = v.NicInfo[0].SecurityGroupType
			tfMap["security_group_ids"] = flex.FlattenStringValueSet(v.NicInfo[0].SecurityGroupIds)
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func encryptPassword(d *schema.ResourceData, client *hpas.Client) string {
	password := d.Get("password").(string)
	secretKey := client.Config.Credentials.SecretAccessKey
	encryptedPassword, _ := bcc.Aes128EncryptUseSecreteKey(secretKey, password)

	return encryptedPassword
}
