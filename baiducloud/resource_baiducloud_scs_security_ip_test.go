package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccScsSecurityIpResourceType = "baiducloud_scs_security_ip"
	testAccScsSecurityIpResourceName = testAccScsSecurityIpResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudScsSecurityIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccScsSecurityIpConfig(BaiduCloudTestResourceTypeNameScsSecurityIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScsSecurityIpResourceName),
					resource.TestCheckResourceAttr(testAccScsSecurityIpResourceName, "instance_id", "scs-BIFDrIl9"),
				),
			},
		},
	})
}

func TestAccBaiduCloudScsSecurityNilIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccScsSecurityNilIpConfig(BaiduCloudTestResourceTypeNameScsSecurityIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScsSecurityIpResourceName),
					resource.TestCheckResourceAttrSet(testAccScsSecurityIpResourceName, "instance_id"),
				),
			},
		},
	})
}

func TestAccBaiduCloudScsSecurityMultiIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccScsSecurityMultiIpConfig(BaiduCloudTestResourceTypeNameScsSecurityIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScsSecurityIpResourceName),
					resource.TestCheckResourceAttrSet(testAccScsSecurityIpResourceName, "instance_id"),
				),
			},
		},
	})
}

func testAccScsSecurityIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_scs_security_ip" "default" {
    instance_id = "scs-bj-hzsywuljybfy"
    security_ips = ["192.168.3.5"]
}
`, name)
}

func testAccScsSecurityNilIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_scs_security_ip" "default" {
    instance_id = "scs-bj-hzsywuljybfy"
}
`, name)
}

func testAccScsSecurityMultiIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_scs_security_ip" "default" {
    instance_id = "scs-bj-hzsywuljybfy"
    security_ips = ["192.168.3.5","192.168.3.6","192.168.3.7"]
}
`, name)
}

func testAccScsSecurityMultiIpFullConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_scs" "default" {
    instance_name           = "%s"
	billing = {
    	payment_timing 		= "Postpaid"
  	}
    purchase_count 			= 1
  	port 					= 6379
	engine_version 			= "3.2"
	node_type 				= "cache.n1.micro"
	cluster_type 			= "master_slave"
	replication_num 		= 1
	shard_num 				= 2
	proxy_num 				= 0
}

resource "baiducloud_scs_security_ip" "default" {
 	instance_id = baiducloud_scs.default.instance_id
    security_ips = ["192.168.3.5","192.168.3.6","192.168.3.7"]
}
`, name)
}

func testAccScsSecurityNilIpFullConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_scs" "default" {
    instance_name           = "%s"
	billing = {
    	payment_timing 		= "Postpaid"
  	}
    purchase_count 			= 1
  	port 					= 6379
	engine_version 			= "3.2"
	node_type 				= "cache.n1.micro"
	cluster_type 			= "master_slave"
	replication_num 		= 1
	shard_num 				= 2
	proxy_num 				= 0
}

resource "baiducloud_scs_security_ip" "default" {
 	instance_id = baiducloud_scs.default.instance_id
}
`, name)
}
