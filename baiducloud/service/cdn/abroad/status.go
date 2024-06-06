package abroad

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
)

func statusAbroadCDNDomain(conn *connectivity.BaiduClient, domainName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		domainConfig, err := FindAbroadDomainConfigByName(conn, domainName)
		if err != nil {
			log.Printf("[ERROR] fail to get abraod CDN domain status: %v", err)
			return nil, "UNKNOWN", err
		}
		return domainConfig, domainConfig.Status, nil
	}
}
