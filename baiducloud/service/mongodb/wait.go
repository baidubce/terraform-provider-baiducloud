package mongodb

import (
	"time"

	"github.com/baidubce/bce-sdk-go/services/mongodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	InstanceAvailableTimeout = 10 * time.Minute
)

func pendingStatus() []string {
	return []string{
		InstanceStatusCreating, InstanceStatusRestarting, InstanceStatusClassChanging,
		InstanceStatusNodeCreating, InstanceStatusNodeRestarting, InstanceStatusNodeClassChanging, InstanceStatusBackuping,
	}
}

func waitInstanceAvailable(conn *connectivity.BaiduClient, instanceID string) (*mongodb.InstanceDetail, error) {
	stateConf := &resource.StateChangeConf{
		Delay:   0,
		Pending: pendingStatus(),
		Target:  []string{InstanceStatusRunning},
		Refresh: statusInstance(conn, instanceID),
		Timeout: InstanceAvailableTimeout,
	}
	raw, err := stateConf.WaitForState()
	if v, ok := raw.(*mongodb.InstanceDetail); ok {
		return v, nil
	}
	return nil, err
}
