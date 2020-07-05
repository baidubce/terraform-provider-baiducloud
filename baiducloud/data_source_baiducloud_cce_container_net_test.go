package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCceContainerNetDataSourceName = "data.baiducloud_cce_container_net.default"
)

//lintignore:AT003
func TestAccBaiduCloudCceContainerNetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNetDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceContainerNetDataSourceName),
					resource.TestCheckResourceAttrSet(testAccCceContainerNetDataSourceName, "container_net"),
					resource.TestCheckResourceAttrSet(testAccCceContainerNetDataSourceName, "capacity"),
				),
			},
		},
	})
}

func testAccContainerNetDataSourceConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "default" {
  name        = "test-BaiduAccVPC"
  description = "created by terraform"
  cidr = "192.168.0.0/16"
}

data "baiducloud_cce_container_net" "default" { 
    vpc_id = baiducloud_vpc.default.id
    vpc_cidr = "192.168.0.0/16"
}
`)
}
