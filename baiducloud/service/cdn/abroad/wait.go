package abroad

import (
	"github.com/baidubce/bce-sdk-go/services/cdn/abroad/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

const (
	CDNDomainAvailableTimeout = 10 * time.Minute
)

func pendingStatus() []string {
	return []string{
		DomainStatusOperating,
	}
}

func waitAbroadCDNDomainAvailable(conn *connectivity.BaiduClient, domainName string) (*api.DomainConfig, error) {
	stateConf := &resource.StateChangeConf{
		Delay:   0,
		Pending: pendingStatus(),
		Target:  []string{DomainStatusRunning,},
		Refresh: statusAbroadCDNDomain(conn, domainName),
		Timeout: CDNDomainAvailableTimeout,
	}
	raw, err := stateConf.WaitForState()
	if v, ok := raw.(*api.DomainConfig); ok {
		return v, nil
	}
	return nil, err
}
