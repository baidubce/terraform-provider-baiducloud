package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccEniAttachmentResourceType = "baiducloud_eni_attachment"
	testAccEniAttachmentResourceName = testAccEniAttachmentResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudEniAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEniAttachmentConfig(BaiduCloudTestResourceTypeNameEniAttachment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEniAttachmentResourceName),
				),
			},
		},
	})
}

func testAccEniAttachmentConfig(name string) string {
	return fmt.Sprintf(`
data "baiducloud_images" "images" {
  image_type = "System"
  name_regex = "8.4 aarch"
  os_name    = "CentOS"
}
resource "baiducloud_vpc" "vpc" {
  name = "terraform_vpc"
  cidr = "172.16.0.0/20"
}
resource "baiducloud_subnet" "subnet" {
  name        = "terraform_subnet"
  zone_name   = "cn-bj-d"
  cidr        = "172.16.0.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}
resource "baiducloud_security_group" "sg" {
  name        = "terraform-sg"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.vpc.id
}
resource "baiducloud_security_group_rule" "sgr1_in" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "ingress"
  source_ip         = "all"
}
resource "baiducloud_security_group_rule" "sgr1_out" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "egress"
  dest_ip           = "all"
}

resource "baiducloud_instance" "server1" {
  availability_zone = "cn-bj-d"
  instance_spec     = "bcc.gr1.c1m4"
  image_id          = data.baiducloud_images.images.images.0.id
  billing           = {
    payment_timing = "Postpaid"
  }
  admin_pass      = "Eni12345"
  subnet_id       = baiducloud_subnet.subnet.id
  security_groups = [
    baiducloud_security_group.sg.id
  ]
}
resource "baiducloud_eni" "eni" {
  name      = "terraform-eni"
  subnet_id = baiducloud_subnet.subnet.id
  #  instance_id = baiducloud_instance.server1.id

  description        = "terraform test"
  security_group_ids = [
    baiducloud_security_group.sg.id
  ]
  private_ip {
    primary            = true
    private_ip_address = ""
  }
}
resource "baiducloud_eni_attachment" "%s" {
  eni_id      = baiducloud_eni.eni.id
  instance_id = baiducloud_instance.server1.id
}
`, name)
}
