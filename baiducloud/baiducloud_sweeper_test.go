package baiducloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClientForRegion returns a common BaiduClient setup needed for the sweeper
// functions for a given region
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("BAIDUCLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty BAIDUCLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("BAIDUCLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty BAIDUCLOUD_SECRET_KEY")
	}

	conf := connectivity.Config{
		Region:    connectivity.Region(region),
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}
