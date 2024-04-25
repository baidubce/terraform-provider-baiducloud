package mongodb

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func statusInstance(conn *connectivity.BaiduClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := findInstance(conn, instanceID)
		if err != nil {
			log.Printf("[ERROR] fail to get instance status: %v", err)
			return nil, "UNKNOWN", err
		}
		return instance, instance.DbInstanceStatus, nil
	}
}
