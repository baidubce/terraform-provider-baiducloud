---
layout: "baiducloud"
subcategory: "ENI"
page_title: "BaiduCloud: baiducloud_eni_attachment"
sidebar_current: "docs-baiducloud-resource-eni_attachment"
description: |-
  Provide a resource to create an ENI association, bind an ENI with instance.
---

# baiducloud_eni_attachment

Provide a resource to create an ENI association, bind an ENI with instance.

~> **NOTE:**
Mount the ENI to the specified cloud host.
A cloud host can be bound to multiple elastic NICs, but can only be bound to one main NIC.
An ENI can only be bound to one cloud host at the same time.
Only a running or powered-off cloud host can bind an ENI.
The elastic network adapter and the bound cloud host must be in the same private network, and the subnets where they are located have the same availability zone.

~> **NOTE:** Depending on the images used by your instance, some images need to be manually configured with an elastic network card to enable the elastic network card bound to the instance to be recognized by the system.
Please refer to: https://cloud.baidu.com/doc/BCC/s/akaz2ccxy

## Example Usage

```hcl
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
resource "baiducloud_security_group_rule" "sgr_in" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "ingress"
}
resource "baiducloud_security_group_rule" "sgr_out" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "all"
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
resource "baiducloud_eip" "eip1" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}
resource "baiducloud_eip" "eip2" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}

resource "baiducloud_eni" "eni" {
  name      = "terraform-eni"
  subnet_id = baiducloud_subnet.subnet.id

  description        = "terraform test"
  security_group_ids = [
    baiducloud_security_group.sg.id
  ]
  private_ip {
    primary            = true
    private_ip_address = "172.16.0.10"
    public_ip_address  = baiducloud_eip.eip2.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.11"
    public_ip_address  = baiducloud_eip.eip1.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.13"
    #    public_ip_address  = baiducloud_eip.eip2.eip
  }
}
resource "time_sleep" "wait_60_seconds" {
  depends_on      = [baiducloud_instance.server1, baiducloud_eni.eni]
  create_duration = "60s"
}
# Wait 60s for the instance to start up completely
resource "baiducloud_eni_attachment" "default" {
  depends_on  = [time_sleep.wait_60_seconds]
  eni_id      = baiducloud_eni.eni.id
  instance_id = baiducloud_instance.server1.id
}
```

## Argument Reference

The following arguments are supported:

* `eni_id` - (Required, ForceNew) Eni ID
* `instance_id` - (Required, ForceNew) Instance ID


