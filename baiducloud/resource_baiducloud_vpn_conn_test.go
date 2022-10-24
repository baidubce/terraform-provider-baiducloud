package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccVPNConnResourceType = "baiducloud_vpn_conn"
	testAccVPNConnResourceName = testAccVPNConnResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudVPNConn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnConfig(BaiduCloudTestResourceTypeNameVPNConn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccVPNConnResourceName),
					resource.TestCheckResourceAttr(testAccVPNConnResourceName, "description", "just for test"),
					resource.TestCheckResourceAttr(testAccVPNConnResourceName, "secret_key", "ddd22@www"),
				),
			},
			{
				Config: testAccVPNConnUpdateConfig(BaiduCloudTestResourceTypeNameVPNConn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccVPNConnResourceName),
					resource.TestCheckResourceAttr(testAccVPNConnResourceName, "description", "just for test new"),
					resource.TestCheckResourceAttr(testAccVPNConnResourceName, "secret_key", "ddd22@qqq"),
				),
			},
		},
	})
}

func testAccVPNConnConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_eip" "default" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}
resource "baiducloud_vpn_gateway" "default" {
  vpn_name       = "test_vpn_gateway"
  vpc_id         = "vpc-65cz3hu92kz2"
  description    = "test desc"
  payment_timing = "Postpaid"
  eip            = baiducloud_eip.default.eip
}
resource "baiducloud_vpn_conn" "default" {
  vpn_id = baiducloud_vpn_gateway.default.id
  secret_key = "ddd22@www"
  local_subnets = ["192.168.0.0/20"]
  remote_ip = "11.11.11.112"
  remote_subnets = ["192.168.100.0/24"]
  description = "just for test"
  vpn_conn_name = "%s"
  ike_config {
    ike_version = "v1"
    ike_mode = "main"
    ike_enc_alg = "aes"
    ike_auth_alg = "sha1"
    ike_pfs = "group2"
    ike_life_time = 100
  }
  ipsec_config {
    ipsec_enc_alg = "aes"
    ipsec_auth_alg = "sha1"
    ipsec_pfs = "group2"
    ipsec_life_time = 200
  }
}
`, name)
}

func testAccVPNConnUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_eip" "default" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}
resource "baiducloud_vpn_gateway" "default" {
  vpn_name       = "test_vpn_gateway"
  vpc_id         = "vpc-65cz3hu92kz2"
  description    = "test desc"
  payment_timing = "Postpaid"
  eip            = baiducloud_eip.default.eip
}
resource "baiducloud_vpn_conn" "default" {
  vpn_id = baiducloud_vpn_gateway.default.id
  secret_key = "ddd22@qqq"
  local_subnets = ["192.168.0.0/20"]
  remote_ip = "11.11.11.112"
  remote_subnets = ["192.168.100.0/24"]
  description = "just for test new"
  vpn_conn_name = "%s"
  ike_config {
    ike_version = "v1"
    ike_mode = "main"
    ike_enc_alg = "aes"
    ike_auth_alg = "sha1"
    ike_pfs = "group2"
    ike_life_time = 100
  }
  ipsec_config {
    ipsec_enc_alg = "aes"
    ipsec_auth_alg = "sha1"
    ipsec_pfs = "group2"
    ipsec_life_time = 200
  }
}
`, name)
}
