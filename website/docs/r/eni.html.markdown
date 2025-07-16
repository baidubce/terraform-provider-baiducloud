---
layout: "baiducloud"
subcategory: "Elastic Network Interface (ENI)"
page_title: "BaiduCloud: baiducloud_eni"
sidebar_current: "docs-baiducloud-resource-eni"
description: |-
  Provide a resource to create an ENI.
---

# baiducloud_eni

Provide a resource to create an ENI.

## Example Usage

```hcl
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
    public_ip_address  = baiducloud_eip.eip1.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.11"
    public_ip_address  = baiducloud_eip.eip2.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.13"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the ENI. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `private_ip` - (Required) Specified intranet IP information
* `subnet_id` - (Required) Subnet ID which ENI belong to
* `description` - (Optional) Description of the ENI
* `enterprise_security_group_ids` - (Optional) Specifies the set of bound enterprise security group IDs
* `security_group_ids` - (Optional) Specifies the set of bound security group IDs

The `private_ip` object supports the following:

* `primary` - (Required) Whether this is the primary private IP address. If true, the IP address cannot be modified. Only one primary IP is allowed per ENI.
* `private_ip_address` - (Optional) The private IP address to assign to the ENI. If empty, one will be assigned automatically.
* `public_ip_address` - (Optional) The public IP address (EIP) of the ENI. Cannot be assigned during creation, or before a private IP is assigned; assign it after the ENI has been created and a private IP is available.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_time` - ENI create time
* `instance_id` - Instance ID the ENI bind
* `mac_address` - Mac address of the ENI
* `status` - Status of ENI, may be inuse, binding, unbinding, available
* `vpc_id` - VPC id which the ENI belong to
* `zone_name` - Availability zone name which ENI belong to


## Import

ENI can be imported, e.g.

```hcl
$ terraform import baiducloud_eni.default eni_id
```

