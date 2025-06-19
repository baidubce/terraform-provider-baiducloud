package hpas

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func StatusInstance(conn *connectivity.BaiduClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := FindInstance(conn, instanceID)
		if err != nil {
			return nil, "UNKNOWN", err
		}
		return instance, instance.Status, nil
	}
}
