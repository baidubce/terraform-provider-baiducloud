package eip

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func StatusEipGroup(conn *connectivity.BaiduClient, eipGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		eipGroup, err := FindEipGroup(conn, eipGroupID)
		if err != nil {
			return nil, "", err
		}
		if eipGroup == nil {
			return nil, "", nil
		}
		return eipGroup, eipGroup.Status, nil
	}
}
