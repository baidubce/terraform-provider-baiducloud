package snic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func StatusSNIC(conn *connectivity.BaiduClient, snicId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snic, err := FindSNIC(conn, snicId)
		if err != nil {
			return nil, "Unknown", err
		}
		return snic, snic.Status, nil
	}
}
