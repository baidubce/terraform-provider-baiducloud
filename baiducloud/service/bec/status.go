package bec

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func StatusVMInstance(conn *connectivity.BaiduClient, vmID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vmInstance, err := FindVMInstance(conn, vmID)
		if err != nil {
			return nil, "UNKNOWN", err
		}
		return vmInstance, vmInstance.Status, nil
	}
}
